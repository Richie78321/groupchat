package memory

import (
	"sync"

	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/Richie78321/groupchat/server/chatdata"
	"github.com/google/uuid"
)

type memoryChatroom struct {
	roomName      string
	lock          sync.Mutex
	subscriptions map[uuid.UUID]chatdata.Subscription

	// Messages in chronological order
	messages       []chatdata.Message
	messagesByUuid map[uuid.UUID]chatdata.Message
}

func newMemoryChatroom(roomName string) *memoryChatroom {
	return &memoryChatroom{
		roomName:      roomName,
		lock:          sync.Mutex{},
		subscriptions: make(map[uuid.UUID]chatdata.Subscription),

		messages:       make([]chatdata.Message, 0),
		messagesByUuid: make(map[uuid.UUID]chatdata.Message),
	}
}

func (c *memoryChatroom) RoomName() string {
	return c.roomName
}

func (c *memoryChatroom) GetLock() sync.Locker {
	return &c.lock
}

func (c *memoryChatroom) SignalSubscriptions() {
	for _, subscription := range c.subscriptions {
		subscription.SignalUpdate()
	}
}

func (c *memoryChatroom) AddSubscription(s chatdata.Subscription) {
	c.subscriptions[s.Id()] = s
}

func (c *memoryChatroom) RemoveSubscription(u uuid.UUID) {
	delete(c.subscriptions, u)
}

func (c *memoryChatroom) Users() (users []*pb.User) {
	for _, subscription := range c.subscriptions {
		users = append(users, subscription.User())
	}

	return users
}

func (c *memoryChatroom) AppendMessage(author *pb.User, body string) error {
	newMessage := newMemoryMessage(author, body)

	c.messages = append(c.messages, newMessage)
	c.messagesByUuid[newMessage.Id()] = newMessage

	return nil
}

func (c *memoryChatroom) LatestMessages(n int) ([]chatdata.Message, error) {
	index := len(c.messages) - n
	if index <= 0 {
		return c.messages, nil
	}

	return c.messages[index:], nil
}

func (c *memoryChatroom) AllMessages() ([]chatdata.Message, error) {
	return c.messages, nil
}

func (c *memoryChatroom) MessageById(u uuid.UUID) (chatdata.Message, bool, error) {
	m, ok := c.messagesByUuid[u]
	return m, ok, nil
}
