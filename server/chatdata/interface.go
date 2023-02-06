package chatdata

import (
	"sync"

	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/google/uuid"
)

type Manager interface {
	GetLock() sync.Locker

	GetRoom(roomName string) (Chatroom, bool)
	CreateRoom(roomName string) Chatroom
}

type Chatroom interface {
	GetLock() sync.Locker

	SignalSubscriptions()
	AddSubscription(Subscription)
	RemoveSubscription(uuid.UUID)
	GetUsers() []*pb.User
}

type Subscription interface {
	Id() uuid.UUID
	User() *pb.User

	SignalUpdate()
	ShouldUpdate() <-chan struct{}
}
