package sqlite

import (
	"sync"

	"github.com/Richie78321/groupchat/server/chatdata"
	"github.com/Richie78321/groupchat/server/chatdata/ephemeralstate"
)

type chatdataManager struct {
	lock sync.Mutex

	sqlChatdata *SqliteChatdata
	esManager   *ephemeralstate.ESManager

	chatrooms map[string]*chatroom
}

func NewChatdataManager(chatdata *SqliteChatdata, esManager *ephemeralstate.ESManager) chatdata.Manager {
	return &chatdataManager{
		lock: sync.Mutex{},

		sqlChatdata: chatdata,
		esManager:   esManager,

		chatrooms: make(map[string]*chatroom),
	}
}

func (m *chatdataManager) GetLock() sync.Locker {
	return &m.lock
}

func (m *chatdataManager) CreateRoom(roomName string) chatdata.Chatroom {
	chatroom := newChatroom(m.sqlChatdata, m.esManager, roomName)
	m.chatrooms[roomName] = chatroom

	return chatroom
}

func (m *chatdataManager) Room(roomName string) (chatdata.Chatroom, bool) {
	room, ok := m.chatrooms[roomName]
	return room, ok
}

func (m *chatdataManager) SignalSubscriptions(roomName string) {
	chatroom, ok := m.Room(roomName)
	if !ok {
		return
	}

	chatroom.SignalSubscriptions()
}
