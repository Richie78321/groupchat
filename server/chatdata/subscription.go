package chatdata

import (
	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/Richie78321/groupchat/server/util"
	"github.com/google/uuid"
)

type subscription struct {
	id   uuid.UUID
	user *pb.User

	update *util.Signal
}

func NewSubscription(user *pb.User) Subscription {
	return &subscription{
		id:     uuid.New(),
		user:   user,
		update: util.NewSignal(),
	}
}

func (s *subscription) Id() uuid.UUID {
	return s.id
}

func (s *subscription) User() *pb.User {
	return s.user
}

func (s *subscription) SignalUpdate() {
	s.update.Signal()
}

func (s *subscription) ShouldUpdate() <-chan struct{} {
	return s.update.GetSignal()
}
