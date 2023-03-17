package sqlite

import (
	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/google/uuid"
)

type message struct {
	sqlChatdata *SqliteChatdata
	event       *MessageEvent
}

func newMessage(event *MessageEvent, sqlChatdata *SqliteChatdata) *message {
	return &message{
		sqlChatdata: sqlChatdata,
		event:       event,
	}
}

func (m *message) Id() uuid.UUID {
	return uuid.MustParse(m.event.MessageID)
}

func (m *message) Author() *pb.User {
	return &pb.User{
		Username: m.event.AuthorID,
	}
}

func (m *message) Body() string {
	return m.event.MessageBody
}

func (m *message) Likers() []*pb.User {
	// TODO(richie): Implement this method
	return nil
}

func (m *message) Like(u *pb.User) bool {
	// TODO(richie): Implement this method
	return false
}

func (m *message) Unlike(u *pb.User) bool {
	// TODO(richie): Implement this method
	return false
}
