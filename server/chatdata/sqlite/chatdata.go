package sqlite

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/Richie78321/groupchat/server/chatdata"
)

type Event struct {
	// Composite primary key out of the PID and sequence number, as this
	// combination should be unique.
	Pid              string `gorm:"primaryKey;autoIncrement:false;not null;index"`
	SequenceNumber   int64  `gorm:"primaryKey;autoIncrement:false;not null;index"`
	LamportTimestamp int64  `gorm:"not null;index"`

	EventType string
	EventID   int
}

type MessageEvent struct {
	ID          int
	ChatroomID  string
	MessageID   string
	AuthorID    string
	MessageBody string
	Event       Event `gorm:"polymorphic:Event"`
}

type LikeEvent struct {
	ID         int
	ChatroomID string
	MessageID  string
	LikerID    string
	Like       bool
	Event      Event `gorm:"polymorphic:Event"`
}

type SqliteChatdata struct {
	lock                 sync.Mutex
	db                   *gorm.DB
	myPid                string
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

	c := &SqliteChatdata{
		lock:                 sync.Mutex{},
		db:                   db,
		myPid:                myPid,
		nextSequenceNumber:   0,
		nextLamportTimestamp: 0,
		incomingEvents:       make(chan *pb.Event, eventBufferSize),
		outgoingEvents:       make(chan *pb.Event, eventBufferSize),
		allPids:              append(otherPids, myPid),
		log:                  log.New(os.Stdout, "[Chatdata] ", log.Default().Flags()),
	}
	if err = c.loadFromDisk(); err != nil {
		return nil, err
	}

	// Spawn a thread to consume new events
	go c.consumeNewEvents()

	return c, nil
}

func (c *SqliteChatdata) IncomingEvents() chan<- *pb.Event {
	return c.incomingEvents
}

func (c *SqliteChatdata) OutgoingEvents() <-chan *pb.Event {
	return c.outgoingEvents
}

func (c *SqliteChatdata) SequenceNumberVector() (chatdata.SequenceNumberVector, error) {
	vector := make(chatdata.SequenceNumberVector)
	for _, pid := range c.allPids {
		nextExpected, err := c.nextExpectedSequenceNumber(pid)
		if err != nil {
			return nil, err
		}

		vector[pid] = nextExpected
	}

	return vector, nil
}

func (c *SqliteChatdata) nextExpectedSequenceNumber(pid string) (int64, error) {
	var result struct {
		NextExpected int64 `gorm:"column:next_expected"`
	}

	// Returns the smallest missing sequence number for a specific PID.
	// This query could likely be made more efficient, but it works for now.
	err := c.db.Raw(`
		SELECT MIN(sequence_number + 1) AS next_expected
		FROM events t1
		WHERE pid = ?
		  AND NOT EXISTS (
			SELECT *
			FROM events t2
			WHERE t2.sequence_number = t1.sequence_number + 1
		)
	`, pid).Scan(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// There are no events with this PID. Therefore the next-expected
			// sequence number is zero.
			return 0, nil
		}

		return 0, err
	}

	return result.NextExpected, nil
}

func (c *SqliteChatdata) EventDiff(vector chatdata.SequenceNumberVector) ([]*pb.Event, error) {
	if vector == nil {
		vector = make(chatdata.SequenceNumberVector)
	}
	// Fill missing PIDs with a next-expected sequence number of 0
	for _, pid := range c.allPids {
		if _, ok := vector[pid]; ok {
			continue
		}
		vector[pid] = 0
	}

	// Get the missing events for every PID
	eventDiff := make([]*pb.Event, 0)
	for pid, sequence_number := range vector {
		events, err := c.sequenceNumberDiff(pid, sequence_number)
		if err != nil {
			return nil, err
		}

		eventDiff = append(eventDiff, events...)
	}

	return eventDiff, nil
}

func (c *SqliteChatdata) sequenceNumberDiff(pid string, sequence_number int64) ([]*pb.Event, error) {
	eventDiff := make([]*pb.Event, 0)

	// Query for MessageAppend events
	messageEvents := make([]MessageEvent, 0)
	err := c.db.Model(&MessageEvent{}).Joins("Event").Where("Event__pid = ? AND Event__sequence_number >= ?", pid, sequence_number).Find(&messageEvents).Error
	if err != nil {
		return nil, err
	}
	for _, messageEvent := range messageEvents {
		eventDiff = append(eventDiff, messageEventToEventPb(&messageEvent))
	}

	// TODO(richie): Finish by also querying for MessageLike events and converting them.
	// Then initialization is complete. Should then make a test to ensure that initial synchronization works.
	//
	// Then need to write a test that shows that updates work.

	return eventDiff, nil
}

func (c *SqliteChatdata) loadFromDisk() error {
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

// consumeNewEvents consumes events and broadcasts new events.
func (c *SqliteChatdata) consumeNewEvents() {
	for {
		newEvent := <-c.incomingEvents
		ignored, err := c.consumeEvent(newEvent)
		if err != nil {
			c.log.Fatalf("%v", err)
		}
		if ignored {
			c.log.Printf("Ignored duplicate event PID=%s, SEQ=%d, LTS=%d", newEvent.Pid, newEvent.SequenceNumber, newEvent.LamportTimestamp)
			continue
		}

		c.log.Printf("Consumed new event PID=%s, SEQ=%d, LTS=%d", newEvent.Pid, newEvent.SequenceNumber, newEvent.LamportTimestamp)

		// If the event was not ignored, then broadcast the event to peers.
		c.outgoingEvents <- newEvent
	}
}

// consumeEvent consumes an event locally. Returns a boolean that is true when
// the event was ignored.
//
// An event is ignored if it has already been consumed, which is represented by
// an event with the same PID and sequence number (the event's composite primary
// key) already existing in the database.
func (c *SqliteChatdata) consumeEvent(event *pb.Event) (bool, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Ensure this event has not already been seen by this server.
	result := c.db.Where("pid = ? AND sequence_number = ?", event.Pid, event.SequenceNumber).Limit(1).Find(&Event{})
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected > 0 {
		// This event has already been seen. Ignore.
		return true, nil
	}

	dbEvent := Event{
		Pid:              event.Pid,
		SequenceNumber:   event.SequenceNumber,
		LamportTimestamp: event.LamportTimestamp,
	}

	// Update this server's Lamport Timestamp if this event has a new maximum timestamp.
	newLts := dbEvent.LamportTimestamp + 1
	if c.nextLamportTimestamp < newLts {
		c.nextLamportTimestamp = newLts
	}

	var convertedEvent interface{}
	switch e := event.Event.(type) {
	case *pb.Event_MessageAppend:
		convertedEvent = &MessageEvent{
			ChatroomID:  e.MessageAppend.ChatroomId,
			MessageID:   e.MessageAppend.MessageUuid,
			AuthorID:    e.MessageAppend.AuthorId,
			MessageBody: e.MessageAppend.Body,
			Event:       dbEvent,
		}
	case *pb.Event_MessageLike:
		convertedEvent = &LikeEvent{
			ChatroomID: e.MessageLike.ChatroomId,
			MessageID:  e.MessageLike.MessageUuid,
			LikerID:    e.MessageLike.LikerId,
			Like:       e.MessageLike.Like,
		}
	default:
		return false, fmt.Errorf("unknown event type: %v", e)
	}

	return false, c.db.Create(convertedEvent).Error
}
