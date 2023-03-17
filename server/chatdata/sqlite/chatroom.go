package sqlite

import (
	"sync"

	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/Richie78321/groupchat/server/chatdata"
	"github.com/google/uuid"
)

type chatroom struct {
	sqlChatdata   *SqliteChatdata
	chatroomId    string
	subscriptions map[uuid.UUID]chatdata.Subscription
}

func newChatroom(sqlChatdata *SqliteChatdata, chatroomId string) *chatroom {
	return &chatroom{
		sqlChatdata:   sqlChatdata,
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
	// TODO(richie): Update ephemeral data when this changes
	c.subscriptions[s.Id()] = s
}

func (c *chatroom) RemoveSubscription(u uuid.UUID) {
	// TODO(richie): Update ephemeral data when this changes
	delete(c.subscriptions, u)
}

func (c *chatroom) Users() (users []*pb.User) {
	// TODO(richie): Implement this using ephemeral data from peers
	return nil
}

func (c *chatroom) AppendMessage(author *pb.User, body string) {
	// TODO(richie): Implement
}

func (c *chatroom) LatestMessages(n int) []chatdata.Message {
	// TODO(richie): Implement
	return nil
}

func (c *chatroom) AllMessages() []chatdata.Message {
	// TODO(richie): Implement
	return nil
}

func (c *chatroom) MessageById(u uuid.UUID) (chatdata.Message, bool) {
	// TODO(richie): Implement
	return nil, false
}
