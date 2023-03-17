package sqlite

import (
	"log"
	"os"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	pb "github.com/Richie78321/groupchat/chatservice"
)

type SqliteChatdata struct {
	// globalLock is held when inserting an event into the database or
	// accessing / mutating the current sequence number or LTS.
	globalLock sync.Mutex
	// chatroomLocks are locks for individual chatrooms. Events for a chatroom
	// should only be consumed if this lock is held.
	chatroomLocks sync.Map

	db    *gorm.DB
	myPid string

	nextSequenceNumber   int64
	nextLamportTimestamp int64

	allPids []string

	incomingEvents chan *pb.Event
	outgoingEvents chan *pb.Event

	log *log.Logger
}

const eventBufferSize = 100

func NewSqliteChatdata(dbPath string, myPid string, otherPids []string) (*SqliteChatdata, error) {
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

	return c, nil
}

// useNextSequenceNumber retrieves and uses the next event sequence number
// and Lamport Timestamp.
func (c *SqliteChatdata) useNextSequenceNumber() (seq int64, lts int64) {
	c.globalLock.Lock()
	defer c.globalLock.Unlock()

	seq = c.nextSequenceNumber
	lts = c.nextLamportTimestamp

	// Increment both values as they have now both been used.
	c.nextSequenceNumber += 1
	c.nextLamportTimestamp += 1

	return seq, lts
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
