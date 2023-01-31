package memory

import (
	"sync"

	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/Richie78321/groupchat/server/chatdata"
	"github.com/google/uuid"
)

type memoryChatroom struct {
	lock          sync.Mutex
	subscriptions map[uuid.UUID]chatdata.Subscription
}

func newMemoryChatroom(roomName string) *memoryChatroom {
	return &memoryChatroom{
		lock:          sync.Mutex{},
		subscriptions: make(map[uuid.UUID]chatdata.Subscription),
	}
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
	// Order matters here: we want to signal all existing subscriptions
	c.SignalSubscriptions()
	c.subscriptions[s.Id()] = s
}

func (c *memoryChatroom) RemoveSubscription(u uuid.UUID) {
	// Order matters here: we want to signal only the remaining subscriptions
	delete(c.subscriptions, u)
	c.SignalSubscriptions()
}

func (c *memoryChatroom) GetUsers() (users []*pb.User) {
	for _, subscription := range c.subscriptions {
		users = append(users, subscription.User())
	}

	return users
}
