package chatserver

import (
	"context"

	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ChatServer) LikeChat(ctx context.Context, req *pb.LikeChatRequest) (*pb.LikeChatResponse, error) {
	// Ensure the message UUID is valid
	messageUuid, err := uuid.Parse(req.MessageUuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "message uuid is invalid: %v", err)
	}

	// Get the chatroom if it exists
	s.manager.GetLock().Lock()
	chatroom, err := s.getChatroomOrFail(req.Chatroom)
	s.manager.GetLock().Unlock()
	if err != nil {
		return nil, err
	}

	chatroom.GetLock().Lock()
	defer chatroom.GetLock().Unlock()

	// Get the message from the chatroom by UUID
	message, ok, err := chatroom.MessageById(messageUuid)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, status.Errorf(codes.NotFound, "could not find message with uuid `%v` in chatroom `%s`", messageUuid, req.Chatroom.Name)
	}

	// Ensure the liker is not the author of the message
	if message.Author().Username == req.Self.Username {
		return nil, status.Error(codes.PermissionDenied, "cannot like your own message")
	}

	var updated bool
	if req.Like {
		updated, err = message.Like(req.Self)
	} else {
		updated, err = message.Unlike(req.Self)
	}
	if err != nil {
		return nil, err
	}

	if updated {
		// Signal the subscriptions because the message has been updated
		chatroom.SignalSubscriptions()
	}

	return &pb.LikeChatResponse{}, nil
}
