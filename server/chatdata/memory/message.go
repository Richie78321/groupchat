package memory

import (
	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/google/uuid"
)

type memoryMessage struct {
	id     uuid.UUID
	author *pb.User
	body   string
}

func newMemoryMessage(author *pb.User, body string) *memoryMessage {
	return &memoryMessage{
		id:     uuid.New(),
		author: author,
		body:   body,
	}
}

func (m *memoryMessage) Id() uuid.UUID {
	return m.id
}

func (m *memoryMessage) Author() *pb.User {
	return m.author
}

func (m *memoryMessage) Body() string {
	return m.body
}
