package sqlite

import (
	"sync"

	"gorm.io/gorm"
)

type SqliteChatroom struct {
	Lock     sync.Mutex
	roomName string
	db       *gorm.DB
}

func newSqliteChatroom(roomName string, db *gorm.DB) *SqliteChatroom {
	return &SqliteChatroom{
		Lock:     sync.Mutex{},
		roomName: roomName,
		db:       db,
	}
}
