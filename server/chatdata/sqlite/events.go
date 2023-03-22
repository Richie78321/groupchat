package sqlite

import (
	"fmt"
	"sync"

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

func pbToEvent(event *pb.Event) (interface{}, error) {
	dbEvent := Event{
		Pid:              event.Pid,
		SequenceNumber:   event.SequenceNumber,
		LamportTimestamp: event.LamportTimestamp,
		ChatroomID:       event.ChatroomId,
	}

	switch e := event.Event.(type) {
	case *pb.Event_MessageAppend:
		return &MessageEvent{
			MessageID:   e.MessageAppend.MessageUuid,
			AuthorID:    e.MessageAppend.AuthorId,
			MessageBody: e.MessageAppend.Body,
			Event:       dbEvent,
		}, nil
	case *pb.Event_MessageLike:
		return &LikeEvent{
			MessageID: e.MessageLike.MessageUuid,
			LikerID:   e.MessageLike.LikerId,
			Like:      e.MessageLike.Like,
		}, nil
	default:
		return nil, fmt.Errorf("unknown event type: %v", e)
	}
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

func (c *SqliteChatdata) ChatroomLock(chatroomId string) sync.Locker {
	chatroomLock, _ := c.chatroomLocks.LoadOrStore(chatroomId, &sync.Mutex{})
	return chatroomLock.(*sync.Mutex)
}

// consumeEvents consumes events from the incomingEvents channel.
func (c *SqliteChatdata) consumeEvents() {
	for {
		newEvent := <-c.incomingEvents

		// Only consume the event when the associated chatroom lock is held.
		// This is to avoid race conditions if multiple threads are attempting to
		// access or modify the chatroom state.
		// The other main user of `c.ConsumeEvent` is the chatserver, which emits
		// new events. The chatserver must also strictly hold the associated chatroom
		// lock when consuming these new events.
		c.ChatroomLock(newEvent.ChatroomId).Lock()
		c.ConsumeEvent(newEvent)
		c.ChatroomLock(newEvent.ChatroomId).Unlock()
	}
}

// ConsumeNewEvent is a helper function to create a new event and consume it.
// This function handles populating event metadata like sequence numbers and LTS.
//
// This function assumes that the caller is currently holding the associated chatroom lock.
func (c *SqliteChatdata) ConsumeNewEvent(event *pb.Event, chatroomId string) error {
	seq, lts := c.useNextSequenceNumber()
	event.ChatroomId = chatroomId
	event.Pid = c.myPid
	event.SequenceNumber = seq
	event.LamportTimestamp = lts

	// `c.ConsumeEvent` is called with the assumption that the associated
	// chatroom lock is already held by the caller of `ConsumeNewEvent`.
	return c.ConsumeEvent(event)
}

// ConsumeEvent is the main method for adding events to the chatdata.
//
// This method is synchronous (meaning that it only returns after the event
// has been consumed). This method also assumes that you've already locked
// the associated chatroom lock for this event.
func (c *SqliteChatdata) ConsumeEvent(event *pb.Event) error {
	ignored, err := c.consumeEventHelper(event)
	if err != nil {
		return err
	}
	if ignored {
		c.log.Printf("Ignored duplicate event PID=%s, SEQ=%d, LTS=%d", event.Pid, event.SequenceNumber, event.LamportTimestamp)
		return nil
	}

	c.log.Printf("Consumed new event PID=%s, SEQ=%d, LTS=%d", event.Pid, event.SequenceNumber, event.LamportTimestamp)

	// Signal chatroom subscriptions to update because a new event was consumed.
	c.SubscriptionSignal.SignalSubscriptions(event.ChatroomId)

	// If the event was not ignored, then broadcast the event to peers.
	c.outgoingEvents <- event
	return nil
}

// consumeEventHelper consumes an event locally. Returns a boolean that is true when
// the event was ignored.
//
// An event is ignored if it has already been consumed, which is represented by
// an event with the same PID and sequence number (the event's composite primary
// key) already existing in the database.
func (c *SqliteChatdata) consumeEventHelper(event *pb.Event) (bool, error) {
	convertedEvent, err := pbToEvent(event)
	if err != nil {
		return false, err
	}

	// Hold the global lock to prevent race conditions with duplicate events
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

	// Update this server's Lamport Timestamp if this event has a new maximum timestamp.
	newLts := event.LamportTimestamp + 1
	if c.nextLamportTimestamp < newLts {
		c.nextLamportTimestamp = newLts
	}

	return false, c.db.Create(convertedEvent).Error
}
