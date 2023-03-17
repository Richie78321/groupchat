package sqlite

import (
	"sync"

	"github.com/Richie78321/groupchat/server/chatdata"
)

type chatdataManager struct {
	lock sync.Mutex

	sqlChatdata *SqliteChatdata
	chatrooms   map[string]*chatroom
}

func NewChatdataManager(chatdata *SqliteChatdata) chatdata.Manager {
	return &chatdataManager{
		lock: sync.Mutex{},

		sqlChatdata: chatdata,
		chatrooms:   make(map[string]*chatroom),
	}
}

func (m *chatdataManager) GetLock() sync.Locker {
	return &m.lock
}

func (m *chatdataManager) CreateRoom(roomName string) chatdata.Chatroom {
	chatroom := newChatroom(m.sqlChatdata, roomName)
	m.chatrooms[roomName] = chatroom

	return chatroom
}

func (m *chatdataManager) Room(roomName string) (chatdata.Chatroom, bool) {
	room, ok := m.chatrooms[roomName]
	return room, ok
}
