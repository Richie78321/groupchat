package sqlite

import "gorm.io/gorm"

type chatroom struct {
	gorm.Model
	ID string
}
