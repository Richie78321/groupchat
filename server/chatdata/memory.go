package chatdata

import (
	"sync"

	"github.com/google/uuid"
)

type memoryManager struct {
	lock      sync.Mutex
	chatrooms map[string]*memoryChatroom
}

func NewMemoryManager() Manager {
	return &memoryManager{
		lock:      sync.Mutex{},
		chatrooms: make(map[string]*memoryChatroom),
	}
}

func (m *memoryManager) GetLock() sync.Locker {
	return &m.lock
}

func (m *memoryManager) GetOrCreateRoom(roomName string) Chatroom {
	if room, ok := m.chatrooms[roomName]; ok {
		return room
	}

	newRoom := newMemoryChatroom(roomName)
	m.chatrooms[roomName] = newRoom
	return newRoom
}

type memoryChatroom struct {
	lock          sync.Mutex
	subscriptions map[uuid.UUID]Subscription
}

func newMemoryChatroom(roomName string) *memoryChatroom {
	return &memoryChatroom{
		lock:          sync.Mutex{},
		subscriptions: make(map[uuid.UUID]Subscription),
	}
}

func (c *memoryChatroom) GetLock() sync.Locker {
	return &c.lock
}

func (c *memoryChatroom) SignalSubscriptions() {
	for _, subscription := range c.subscriptions {
		subscription.SignalUpdate()
	}
}

func (c *memoryChatroom) AddSubscription(s Subscription) {
	c.subscriptions[s.Id()] = s
}

func (c *memoryChatroom) RemoveSubscription(u uuid.UUID) {
	delete(c.subscriptions, u)
}
