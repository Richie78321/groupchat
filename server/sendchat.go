package server

import (
	pb "github.com/Richie78321/groupchat/chatservice"
	"golang.org/x/net/context"
)

func (s *chatServer) SendChat(ctx context.Context, req *pb.SendChatRequest) (*pb.SendChatResponse, error) {
	// Get the chatroom if it exists
	s.manager.GetLock().Lock()
	chatroom, err := s.getChatroomOrFail(req.Chatroom)
	s.manager.GetLock().Unlock()
	if err != nil {
		return nil, err
	}

	chatroom.GetLock().Lock()
	defer chatroom.GetLock().Unlock()

	// Ensure the user is logged in
	err = ensureUserLoggedIn(chatroom, req.Self)
	if err != nil {
		return nil, err
	}

	// Create the message and append it to the chatroom
	chatroom.AppendMessage(req.Self, req.Body)

	// A new message has been added, so signal the subscribers
	chatroom.SignalSubscriptions()

	return &pb.SendChatResponse{}, nil
}
