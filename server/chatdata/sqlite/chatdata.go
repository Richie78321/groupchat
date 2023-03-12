package sqlite

import (
	"fmt"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	pb "github.com/Richie78321/groupchat/chatservice"
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

type SqliteChatdata struct {
	Lock                 sync.Mutex
	db                   *gorm.DB
	pid                  string
	nextSequenceNumber   int64
	nextLamportTimestamp int64
}

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
	}
	if err = c.loadFromDisk(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *SqliteChatdata) loadFromDisk() error {
	// Load the next sequence number from disk.
	selectedEvent := &Event{}
	result := c.db.Order("sequence_number DESC").Where("pid = ?", c.pid).Limit(1).Find(selectedEvent)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected > 0 {
		// Set the next sequence number as the current maximum sequence number plus 1
		c.nextSequenceNumber = selectedEvent.SequenceNumber + 1
	}

	// TODO: Load the Lamport Timestamp

	return nil
}

func (c *SqliteChatdata) ConsumeEvent(event *pb.Event) (bool, error) {
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

	switch e := event.Event.(type) {
	case *pb.Event_MessageAppend:
		return false, c.consumeMessageAppend(&MessageEvent{
			ChatroomID:  e.MessageAppend.ChatroomId,
			MessageID:   e.MessageAppend.MessageUuid,
			AuthorID:    e.MessageAppend.AuthorId,
			MessageBody: e.MessageAppend.Body,
			Event:       dbEvent,
		})
	case *pb.Event_MessageLike:
		// TODO(richie): Implement
		return false, nil
	default:
		return false, fmt.Errorf("unknown event type: %v", e)
	}
}

func (c *SqliteChatdata) consumeMessageAppend(event *MessageEvent) error {
	// Add the message event to the database
	return c.db.Create(event).Error
}
