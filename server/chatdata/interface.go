package chatdata

import (
	"sync"

	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/google/uuid"
)

type Manager interface {
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
	LatestMessages(int) []Message
	AllMessages() []Message
	MessageById(uuid.UUID) (Message, bool)
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

	Likers() []*pb.User
	Like(*pb.User) bool
	Unlike(*pb.User) bool
}
