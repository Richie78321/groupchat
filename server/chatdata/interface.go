package chatdata

import (
	"sync"

	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/google/uuid"
)

type SubscriptionSignal interface {
	SignalSubscriptions(string)
}

type Manager interface {
	SubscriptionSignal
	GetLock() sync.Locker

	Room(roomName string) (Chatroom, bool)
	CreateRoom(roomName string) Chatroom
}

type Chatroom interface {
	GetLock() sync.Locker

	RoomName() string

	SignalSubscriptions()
	AddSubscription(Subscription)
	RemoveSubscription(uuid.UUID)
	Users() []*pb.User

	AppendMessage(*pb.User, string) error
	LatestMessages(int) ([]Message, error)
	AllMessages() ([]Message, error)
	MessageById(uuid.UUID) (Message, bool, error)
}

type Subscription interface {
	Id() uuid.UUID
	User() *pb.User

	SignalUpdate()
	ShouldUpdate() <-chan struct{}
}

type Message interface {
	Id() uuid.UUID
	Author() *pb.User
	Body() string

	Likers() ([]*pb.User, error)
	Like(*pb.User) (bool, error)
	Unlike(*pb.User) (bool, error)
}
