package sqlite

import (
	"fmt"
	"log"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/Richie78321/groupchat/server/chatdata"
)

type Event struct {
	// Composite primary key out of the PID and sequence number, as this
	// combination should be unique.
	Pid              string `gorm:"primaryKey;autoIncrement:false;not null"`
	SequenceNumber   int64  `gorm:"primaryKey;autoIncrement:false;not null"`
	LamportTimestamp int64  `gorm:"not null"`

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
	Lock                 sync.Mutex
	db                   *gorm.DB
	pid                  string
	nextSequenceNumber   int64
	nextLamportTimestamp int64

	incomingEvents chan *pb.Event
	outgoingEvents chan *pb.Event
}

const eventBufferSize = 100

func NewSqliteChatdata(dbPath string, pid string) (*SqliteChatdata, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Create the necessary models if they do not already exist
	db.AutoMigrate(&Event{})
	db.AutoMigrate(&MessageEvent{})

	c := &SqliteChatdata{
		Lock:                 sync.Mutex{},
		db:                   db,
		pid:                  pid,
		nextSequenceNumber:   0,
		nextLamportTimestamp: 0,
		incomingEvents:       make(chan *pb.Event, eventBufferSize),
		outgoingEvents:       make(chan *pb.Event, eventBufferSize),
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
	return c.incomingEvents
}

func (c *SqliteChatdata) SequenceNumberVector() chatdata.SequenceNumberVector {
	// TODO(richie): Implement
	return nil
}

func (c *SqliteChatdata) EventDiff(vector chatdata.SequenceNumberVector) []*pb.Event {
	// TODO(richie): Implement
	return nil
}

func (c *SqliteChatdata) loadFromDisk() error {
	// Load the next sequence number from disk.
	selectedEvent := &Event{}
	// We are only concerned with events that have the server's PID.
	result := c.db.Order("sequence_number DESC").Where("pid = ?", c.pid).Limit(1).Find(selectedEvent)
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
			log.Fatalf("%v", err)
		}
		if ignored {
			continue
		}

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
	c.Lock.Lock()
	defer c.Lock.Unlock()

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
