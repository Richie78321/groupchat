package memory

import (
	"sync"

	"github.com/Richie78321/groupchat/server/chatdata"
)

type memoryManager struct {
	lock      sync.Mutex
	chatrooms map[string]*memoryChatroom
}

func NewMemoryManager() chatdata.Manager {
	return &memoryManager{
		lock:      sync.Mutex{},
		chatrooms: make(map[string]*memoryChatroom),
	}
}

func (m *memoryManager) GetLock() sync.Locker {
	return &m.lock
}

func (m *memoryManager) GetOrCreateRoom(roomName string) chatdata.Chatroom {
	if room, ok := m.chatrooms[roomName]; ok {
		return room
	}

	newRoom := newMemoryChatroom(roomName)
	m.chatrooms[roomName] = newRoom
	return newRoom
}
