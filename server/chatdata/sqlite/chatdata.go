package sqlite

import (
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/Richie78321/groupchat/server/chatdata"
)

const garbageCollectSleep = 1 * time.Minute

type SqliteChatdata struct {
	// globalLock is held when inserting an event into the database,
	// accessing / mutating the current sequence number or LTS, or
	// accessing / mutating event metadata.
	globalLock sync.Mutex
	// chatroomLocks are locks for individual chatrooms. Events for a chatroom
	// should only be consumed if this lock is held.
	chatroomLocks sync.Map

	// stateCaches is a map from chatroom ID to chatroomCache.
	//
	// It is assumed that you hold the associated chatroom lock when
	// accessing / mutating the chatroom cache.
	stateCaches sync.Map

	db    *gorm.DB
	myPid string

	metadataPath  string
	eventMetadata struct {
		// GarbageCollectedTo maps from PID to the maximum sequence number from this
		// PID where garbage collection ran.
		GarbageCollectedTo chatdata.GarbageCollectedToVector `json:"contiguousUpTo"`
	}

	nextSequenceNumber   int64
	nextLamportTimestamp int64

	allPids []string

	incomingEvents chan *pb.Event
	outgoingEvents chan *pb.Event

	SubscriptionSignal chatdata.SubscriptionSignal

	log *log.Logger
}

const eventBufferSize = 100

func NewSqliteChatdata(dbPath string, metadataPath string, myPid string, otherPids []string) (*SqliteChatdata, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Create the necessary models if they do not already exist
	db.AutoMigrate(&Event{})
	db.AutoMigrate(&MessageEvent{})
	db.AutoMigrate(&LikeEvent{})

	c := &SqliteChatdata{
		globalLock: sync.Mutex{},
		db:         db,
		myPid:      myPid,

		metadataPath: metadataPath,

		nextSequenceNumber:   0,
		nextLamportTimestamp: 0,

		incomingEvents: make(chan *pb.Event, eventBufferSize),
		outgoingEvents: make(chan *pb.Event, eventBufferSize),
		allPids:        append(otherPids, myPid),
		log:            log.New(os.Stdout, "[Chatdata] ", log.Default().Flags()),
	}
	if err = c.loadFromDisk(); err != nil {
		return nil, err
	}

	// Spawn a thread to consume new events
	go c.consumeEvents()

	// Ensure the event metadata is loaded from disk
	c.loadEventMetadata()
	// Spawn a thread to periodically garbage collect
	go c.garbageCollectRoutine()

	return c, nil
}

func (c *SqliteChatdata) getChatroomCache(chatroomId string) *chatroomCache {
	cache, _ := c.stateCaches.LoadOrStore(chatroomId, &chatroomCache{})
	return cache.(*chatroomCache)
}

func (c *SqliteChatdata) loadFromDisk() error {
	// Hold the global lock because we are changing the sequence number and LTS
	c.globalLock.Lock()
	defer c.globalLock.Unlock()

	// Load the next sequence number from disk.
	selectedEvent := &Event{}
	// We are only concerned with events that have the server's PID.
	result := c.db.Order("sequence_number DESC").Where("pid = ?", c.myPid).Limit(1).Find(selectedEvent)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected > 0 {
		// Set the next sequence number as the current maximum sequence number plus 1
		c.nextSequenceNumber = selectedEvent.SequenceNumber + 1
	}

	// Load the next Lamport Timestamp from disk.
	selectedEvent = &Event{}
	result = c.db.Order("lamport_timestamp DESC").Limit(1).Find(selectedEvent)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected > 0 {
		// Set the next Lamport Timestamp as the current maximum LTS plus 1
		c.nextLamportTimestamp = selectedEvent.LamportTimestamp + 1
	}

	return nil
}

func reverseList[T any](list []T) {
	for i := 0; i < len(list)/2; i++ {
		list[i], list[len(list)-i-1] = list[len(list)-i-1], list[i]
	}
}

func causalOrdering(tx *gorm.DB) *gorm.DB {
	return tx.
		// Order messages by their LTS to preserve causality and provide a consistent
		// message ordering.
		Order("Event__lamport_timestamp desc").
		// Break ties in LTS with the PID to ensure that message ordering is
		// deterministic across processes.
		Order("Event__pid")
}

func (c *SqliteChatdata) MessageById(chatroomId string, messageId string) (*MessageEvent, error) {
	message := MessageEvent{}
	result := causalOrdering(
		c.db.Model(&MessageEvent{}).
			Joins("Event").
			Where("Event__chatroom_id = ? AND message_id = ?", chatroomId, messageId).
			Limit(1),
	).Find(&message)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected <= 0 {
		return nil, nil
	}

	return &message, nil
}

func (c *SqliteChatdata) GetLatestMessages(chatroomId string, limit int) ([]*MessageEvent, error) {
	latestMessages := make([]*MessageEvent, 0)
	err := causalOrdering(
		c.db.Model(&MessageEvent{}).
			Joins("Event").
			Where("Event__chatroom_id = ?", chatroomId).
			Limit(limit),
	).Find(&latestMessages).Error
	if err != nil {
		return nil, err
	}

	// Reverse the order of the messages so the latest message is the last element of the list
	reverseList(latestMessages)
	return latestMessages, nil
}

func (c *SqliteChatdata) IsLiker(chatroomId string, messageId string, userId string) (bool, error) {
	likeEvent := LikeEvent{}
	result := causalOrdering(
		c.db.Model(&LikeEvent{}).
			Joins("Event").
			// Retrieve like events from this chatroom, for this message, and by this user.
			Where("Event__chatroom_id = ? AND message_id = ? AND liker_id = ?", chatroomId, messageId, userId).
			Limit(1),
	).Find(&likeEvent)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected <= 0 {
		// If there are no like events for this user, that is equivalent to the
		// user not being a liker.
		return false, nil
	}

	return likeEvent.Like, nil
}

func (c *SqliteChatdata) GetLikers(chatroomId string, messageId string) ([]*LikeEvent, error) {
	likeEvents := make([]*LikeEvent, 0)
	err := causalOrdering(
		c.db.Model(&LikeEvent{}).
			Joins("Event").
			Where("Event__chatroom_id = ? AND message_id = ?", chatroomId, messageId).
			Order("liker_id"),
	).Find(&likeEvents).Error
	if err != nil {
		return nil, err
	}

	latestLikeEvents := make([]*LikeEvent, 0)
	for i := 0; i < len(likeEvents); i++ {
		event := likeEvents[i]
		if event.Like {
			// Only retain the positive like events.
			latestLikeEvents = append(latestLikeEvents, event)
		}

		// Keep the first like event from each liker. This works because the
		// like events are returned from the SQL query ordered first by liker
		// and then by causal ordering.
		for i+1 < len(likeEvents) && likeEvents[i+1].LikerID == event.LikerID {
			i++
		}
	}

	return latestLikeEvents, nil
}
