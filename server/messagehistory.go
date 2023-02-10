package server

import (
	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/Richie78321/groupchat/server/chatdata"
	"golang.org/x/net/context"
)

func (s *chatServer) MessageHistory(ctx context.Context, req *pb.MessageHistoryRequest) (*pb.MessageHistoryResponse, error) {
	// Get the chatroom if it exists
	s.manager.GetLock().Lock()
	chatroom, err := s.getChatroomOrFail(req.Chatroom)
	s.manager.GetLock().Unlock()
	if err != nil {
		return nil, err
	}

	chatroom.GetLock().Lock()
	defer chatroom.GetLock().Unlock()

	return &pb.MessageHistoryResponse{
		Messages: chatdata.MessageListToPb(chatroom.AllMessages()),
	}, nil
}
