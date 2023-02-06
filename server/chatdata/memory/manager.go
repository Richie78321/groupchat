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

func (m *memoryManager) CreateRoom(roomName string) chatdata.Chatroom {
	chatroom := newMemoryChatroom(roomName)
	m.chatrooms[roomName] = chatroom

	return chatroom
}

func (m *memoryManager) Room(roomName string) (chatdata.Chatroom, bool) {
	room, ok := m.chatrooms[roomName]
	return room, ok
}
