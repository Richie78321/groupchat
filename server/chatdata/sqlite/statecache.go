package sqlite

import (
	"fmt"
)

type chatroomCache struct {
	messageOrder      []string
	messagesById      map[string]*MessageEvent
	likersByMessageId map[string][]*LikeEvent
}

func (s *chatroomCache) InvalidateCache(newEvent interface{}) error {
	switch event := newEvent.(type) {
	case *MessageEvent:
		// With a new message, the message order will change.
		s.messageOrder = nil
		// Remove the existing cache for the message if it exists.
		delete(s.messagesById, event.MessageID)
	case *LikeEvent:
		// Remove the existing cache for a message's likers if it exists.
		delete(s.likersByMessageId, event.MessageID)
	default:
		return fmt.Errorf("unknown event type %v", event)
	}

	return nil
}
