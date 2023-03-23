package memory

import (
	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/google/uuid"
)

type memoryMessage struct {
	id     uuid.UUID
	author *pb.User
	body   string

	likersByUsername map[string]*pb.User
}

func newMemoryMessage(author *pb.User, body string) *memoryMessage {
	return &memoryMessage{
		id:               uuid.New(),
		author:           author,
		body:             body,
		likersByUsername: make(map[string]*pb.User),
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

func (m *memoryMessage) Likers() ([]*pb.User, error) {
	likers := make([]*pb.User, 0, len(m.likersByUsername))
	for _, user := range m.likersByUsername {
		likers = append(likers, user)
	}

	return likers, nil
}

func (m *memoryMessage) Like(u *pb.User) (bool, error) {
	if _, ok := m.likersByUsername[u.Username]; ok {
		return false, nil
	}

	m.likersByUsername[u.Username] = u
	return true, nil
}

func (m *memoryMessage) Unlike(u *pb.User) (bool, error) {
	if _, ok := m.likersByUsername[u.Username]; !ok {
		return false, nil
	}

	delete(m.likersByUsername, u.Username)
	return true, nil
}
