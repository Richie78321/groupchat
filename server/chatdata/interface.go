package chatdata

import (
	"sync"

	"github.com/google/uuid"
)

type Manager interface {
	GetLock() sync.Locker

	GetOrCreateRoom(roomName string) Chatroom
}

type Chatroom interface {
	GetLock() sync.Locker

	SignalSubscriptions()
	AddSubscription(Subscription)
	RemoveSubscription(uuid.UUID)
}

type Message interface {
}

type Subscription interface {
	Id() uuid.UUID
	Username() string

	SignalUpdate()
	ShouldUpdate() <-chan struct{}
}
