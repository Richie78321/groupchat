package sqlite

import (
	"sync"

	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/Richie78321/groupchat/server/chatdata"
	"github.com/google/uuid"
)

type chatroom struct {
	chatdata      *SqliteChatdata
	chatroomId    string
	subscriptions map[uuid.UUID]chatdata.Subscription
}

func newChatroom(chatdata *SqliteChatdata, chatroomId string) *chatroom {
	return &chatroom{
		chatdata:   chatdata,
		chatroomId: chatroomId,
	}
}

func (c *chatroom) GetLock() sync.Locker {
	return c.chatdata.ChatroomLock(c.chatroomId)
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
