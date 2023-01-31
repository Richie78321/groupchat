package chatdata

import (
	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/google/uuid"
)

type subscription struct {
	id   uuid.UUID
	user *pb.User

	// A channel for signalling with a strict buffer size of 1.
	update chan struct{}
}

func NewSubscription(user *pb.User) Subscription {
	return &subscription{
		id:   uuid.New(),
		user: user,

		// The buffer size is strictly 1 to ensure proper signalling behavior.
		update: make(chan struct{}, 1),
	}
}

func (s *subscription) Id() uuid.UUID {
	return s.id
}

func (s *subscription) User() *pb.User {
	return s.user
}

func (s *subscription) SignalUpdate() {
	select {
	case s.update <- struct{}{}:
	default:
		// The update channel has a strict buffer size of 1.
		// If there is already a signal in the subscription update channel,
		// then we can safely continue because the subscription has already
		// been signalled to update.
	}
}

func (s *subscription) ShouldUpdate() <-chan struct{} {
	return s.update
}
