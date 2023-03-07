package sqlite

import (
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqliteChatdata struct {
	Lock      sync.Mutex
	db        *gorm.DB
	chatrooms map[string]*SqliteChatroom
}

func NewSqliteManager(dbPath string) (*SqliteChatdata, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Create the necessary models if they do not already exist
	db.AutoMigrate(&chatroom{})

	// TODO(richie): Need to load from persistent state? Or could just load on first fetch for a chatroom. Pretty much just locks that need to be loaded. Can load those at first use.

	return &SqliteChatdata{
		Lock:      sync.Mutex{},
		db:        db,
		chatrooms: make(map[string]*SqliteChatroom),
	}, nil
}

func (c *SqliteChatdata) Room(roomName string) (*SqliteChatroom, bool) {
	chatroom, ok := c.chatrooms[roomName]
	return chatroom, ok
}

func (c *SqliteChatdata) CreateRoom(roomName string) *SqliteChatroom {
	// Create a record for the chatroom if it does not already exist
	c.db.FirstOrCreate(&chatroom{
		ID: roomName,
	})

	newChatroom := newSqliteChatroom(roomName, c.db)
	c.chatrooms[roomName] = newChatroom
	return newChatroom
}
