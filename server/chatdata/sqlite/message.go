package sqlite

import (
	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/google/uuid"
)

type message struct {
	sqlChatdata *SqliteChatdata
	chatroomId  string
	event       *MessageEvent
}

func newMessage(event *MessageEvent, chatroomId string, sqlChatdata *SqliteChatdata) *message {
	return &message{
		sqlChatdata: sqlChatdata,
		chatroomId:  chatroomId,
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

func (m *message) Likers() ([]*pb.User, error) {
	likers, err := m.sqlChatdata.GetLikers(m.event.Event.ChatroomID, m.event.MessageID)
	if err != nil {
		return nil, err
	}

	likerUsers := make([]*pb.User, 0, len(likers))
	for _, liker := range likers {
		likerUsers = append(likerUsers, &pb.User{
			Username: liker.LikerID,
		})
	}

	return likerUsers, nil
}

func (m *message) likeMessageHelper(u *pb.User, like bool) (bool, error) {
	isLiker, err := m.sqlChatdata.IsLiker(m.chatroomId, m.event.MessageID, u.Username)
	if err != nil {
		return false, err
	}

	if isLiker == like {
		// The user is already in the correct like state, so this can be a no-op.
		return false, nil
	}

	err = m.sqlChatdata.ConsumeNewEvent(&pb.Event{
		Event: &pb.Event_MessageLike{
			MessageLike: &pb.MessageLike{
				MessageUuid: m.event.MessageID,
				LikerId:     u.Username,
				Like:        like,
			},
		},
	}, m.chatroomId)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (m *message) Like(u *pb.User) (bool, error) {
	return m.likeMessageHelper(u, true)
}

func (m *message) Unlike(u *pb.User) (bool, error) {
	return m.likeMessageHelper(u, false)
}
