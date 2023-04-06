package sqlite

import (
	"sync"

	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/Richie78321/groupchat/server/chatdata"
	"github.com/Richie78321/groupchat/server/chatdata/ephemeralstate"
	"github.com/google/uuid"
)

type chatroom struct {
	sqlChatdata   *SqliteChatdata
	esManager     *ephemeralstate.ESManager
	chatroomId    string
	subscriptions map[uuid.UUID]chatdata.Subscription
}

func newChatroom(sqlChatdata *SqliteChatdata, esManager *ephemeralstate.ESManager, chatroomId string) *chatroom {
	return &chatroom{
		sqlChatdata:   sqlChatdata,
		esManager:     esManager,
		chatroomId:    chatroomId,
		subscriptions: make(map[uuid.UUID]chatdata.Subscription),
	}
}

func (c *chatroom) GetLock() sync.Locker {
	return c.sqlChatdata.ChatroomLock(c.chatroomId)
}

func (c *chatroom) RoomName() string {
	return c.chatroomId
}

func (c *chatroom) SignalSubscriptions() {
	for _, subscription := range c.subscriptions {
		subscription.SignalUpdate()
	}
}

func (c *chatroom) AddSubscription(s chatdata.Subscription) {
	c.subscriptions[s.Id()] = s

	// Update the ephemeral state with the new client connection
	c.esManager.ClientConnected(true, s.User())
}

func (c *chatroom) RemoveSubscription(u uuid.UUID) {
	subscription, ok := c.subscriptions[u]
	if !ok {
		return
	}

	delete(c.subscriptions, u)

	// Update the ephemeral state with the client disconnection
	c.esManager.ClientConnected(false, subscription.User())
}

func (c *chatroom) Users() (users []*pb.User) {
	c.esManager.Lock.RLock()
	defer c.esManager.Lock.RUnlock()

	esGroup := c.esManager.ES()
	// Collect the users from each server by username to de-duplicate users
	// that are logged in on multiple servers.
	usersByUsername := make(map[string]*pb.User)
	for _, es := range esGroup {
		if es == nil {
			continue
		}

		for _, user := range es.ConnectedClients {
			usersByUsername[user.Username] = user
		}
	}

	for _, user := range usersByUsername {
		users = append(users, user)
	}

	return users
}

func (c *chatroom) AppendMessage(author *pb.User, body string) error {
	return c.sqlChatdata.ConsumeNewEvent(&pb.Event{
		Event: &pb.Event_MessageAppend{
			MessageAppend: &pb.MessageAppend{
				MessageUuid: uuid.NewString(),
				AuthorId:    author.Username,
				Body:        body,
			},
		},
	}, c.chatroomId)
}

func (c *chatroom) LatestMessages(n int) ([]chatdata.Message, error) {
	messageEvents, err := c.sqlChatdata.LatestMessages(c.chatroomId, n)
	if err != nil {
		return nil, err
	}

	messages := make([]chatdata.Message, len(messageEvents))
	for i, messageEvent := range messageEvents {
		messages[i] = newMessage(messageEvent, c.chatroomId, c.sqlChatdata)
	}

	return messages, nil
}

func (c *chatroom) AllMessages() ([]chatdata.Message, error) {
	// Specify no limit (-1) to get the entire message history.
	return c.LatestMessages(-1)
}

func (c *chatroom) MessageById(u uuid.UUID) (chatdata.Message, bool, error) {
	messageEvent, err := c.sqlChatdata.MessageById(c.chatroomId, u.String())
	if err != nil {
		return nil, false, err
	}
	if messageEvent == nil {
		return nil, false, nil
	}

	return newMessage(messageEvent, c.chatroomId, c.sqlChatdata), true, nil
}
