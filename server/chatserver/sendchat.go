package chatserver

import (
	pb "github.com/Richie78321/groupchat/chatservice"
	"golang.org/x/net/context"
)

func (s *ChatServer) SendChat(ctx context.Context, req *pb.SendChatRequest) (*pb.SendChatResponse, error) {
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
	if err := ensureUserLoggedIn(chatroom, req.Self); err != nil {
		return nil, err
	}

	// Create the message and append it to the chatroom
	if err := chatroom.AppendMessage(req.Self, req.Body); err != nil {
		return nil, err
	}

	// A new message has been added, so signal the subscribers
	chatroom.SignalSubscriptions()

	return &pb.SendChatResponse{}, nil
}
