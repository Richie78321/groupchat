package chatserver

import (
	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/Richie78321/groupchat/server/chatdata"
	"golang.org/x/net/context"
)

func (s *ChatServer) MessageHistory(ctx context.Context, req *pb.MessageHistoryRequest) (*pb.MessageHistoryResponse, error) {
	// Get the chatroom if it exists
	s.manager.GetLock().Lock()
	chatroom, err := s.getChatroomOrFail(req.Chatroom)
	s.manager.GetLock().Unlock()
	if err != nil {
		return nil, err
	}

	chatroom.GetLock().Lock()
	defer chatroom.GetLock().Unlock()

	messages, err := chatroom.AllMessages()
	if err != nil {
		return nil, err
	}

	messagesPb, err := chatdata.MessageListToPb(messages)
	if err != nil {
		return nil, err
	}

	return &pb.MessageHistoryResponse{
		Messages: messagesPb,
	}, nil
}
