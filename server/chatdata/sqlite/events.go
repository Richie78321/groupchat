package sqlite

import (
	"fmt"

	pb "github.com/Richie78321/groupchat/chatservice"
)

type Event struct {
	// Composite primary key out of the PID and sequence number, as this
	// combination should be unique.
	Pid              string `gorm:"primaryKey;autoIncrement:false;not null;index"`
	SequenceNumber   int64  `gorm:"primaryKey;autoIncrement:false;not null;index"`
	LamportTimestamp int64  `gorm:"not null;index"`
	ChatroomID       string

	EventType string
	EventID   int
}

type MessageEvent struct {
	ID          int
	MessageID   string
	AuthorID    string
	MessageBody string
	Event       Event `gorm:"polymorphic:Event"`
}

type LikeEvent struct {
	ID        int
	MessageID string
	LikerID   string
	Like      bool
	Event     Event `gorm:"polymorphic:Event"`
}

func messageEventToPb(m *MessageEvent) *pb.Event {
	return &pb.Event{
		Pid:              m.Event.Pid,
		SequenceNumber:   m.Event.SequenceNumber,
		LamportTimestamp: m.Event.LamportTimestamp,
		ChatroomId:       m.Event.ChatroomID,
		Event: &pb.Event_MessageAppend{
			MessageAppend: &pb.MessageAppend{
				MessageUuid: m.MessageID,
				AuthorId:    m.AuthorID,
				Body:        m.MessageBody,
			},
		},
	}
}

func likeEventToPb(l *LikeEvent) *pb.Event {
	return &pb.Event{
		Pid:              l.Event.Pid,
		SequenceNumber:   l.Event.SequenceNumber,
		LamportTimestamp: l.Event.LamportTimestamp,
		ChatroomId:       l.Event.ChatroomID,
		Event: &pb.Event_MessageLike{
			MessageLike: &pb.MessageLike{
				MessageUuid: l.MessageID,
				LikerId:     l.LikerID,
				Like:        l.Like,
			},
		},
	}
}

// consumeEvents consumes events and broadcasts new events.
func (c *SqliteChatdata) consumeEvents() {
	for {
		newEvent := <-c.incomingEvents
		ignored, err := c.consumeEventHelper(newEvent)
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

// consumeEventHelper consumes an event locally. Returns a boolean that is true when
// the event was ignored.
//
// An event is ignored if it has already been consumed, which is represented by
// an event with the same PID and sequence number (the event's composite primary
// key) already existing in the database.
func (c *SqliteChatdata) consumeEventHelper(event *pb.Event) (bool, error) {
	// Hold the global lock to prevent any race conditions with duplicate events
	// and changes to LTS.
	c.globalLock.Lock()
	defer c.globalLock.Unlock()

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
		ChatroomID:       event.ChatroomId,
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
			MessageID:   e.MessageAppend.MessageUuid,
			AuthorID:    e.MessageAppend.AuthorId,
			MessageBody: e.MessageAppend.Body,
			Event:       dbEvent,
		}
	case *pb.Event_MessageLike:
		convertedEvent = &LikeEvent{
			MessageID: e.MessageLike.MessageUuid,
			LikerID:   e.MessageLike.LikerId,
			Like:      e.MessageLike.Like,
		}
	default:
		return false, fmt.Errorf("unknown event type: %v", e)
	}

	return false, c.db.Create(convertedEvent).Error
}
