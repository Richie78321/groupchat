package server

import (
	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/Richie78321/groupchat/server/chatdata"
	"golang.org/x/net/context"
)

func sendChatHelper(ctx context.Context, c chatdata.Chatroom, req *pb.SendChatRequest) (*pb.SendChatResponse, error) {
	c.GetLock().Lock()
	// Helper created to use defer on unlock.
	defer c.GetLock().Unlock()

	// Ensure the user is logged in
	err := ensureUserLoggedIn(c, req.Self)
	if err != nil {
		return nil, err
	}

	// Create the message and append it to the chatroom
	c.AppendMessage(req.Self, req.Body)

	return &pb.SendChatResponse{}, nil
}

func (s *chatServer) SendChat(ctx context.Context, req *pb.SendChatRequest) (*pb.SendChatResponse, error) {
	// Get the chatroom if it exists
	s.manager.GetLock().Lock()
	chatroom, err := s.getChatroomOrFail(req.Chatroom)
	s.manager.GetLock().Unlock()
	if err != nil {
		return nil, err
	}

	return sendChatHelper(ctx, chatroom, req)
}
