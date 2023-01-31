package server

import (
	"fmt"
	"net"
	"time"

	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/Richie78321/groupchat/server/chatdata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type chatServer struct {
	manager chatdata.Manager
	pb.UnimplementedChatServiceServer
}

func newChatServer() *chatServer {
	return &chatServer{
		manager: chatdata.NewMemoryManager(),
	}
}

func (s *chatServer) sendSubscriptionUpdate(stream pb.ChatService_SubscribeChatroomServer) error {
	return stream.Send(&pb.ChatroomSubscriptionUpdate{})
}

func (s *chatServer) SubscribeChatroom(req *pb.SubscribeChatroomRequest, stream pb.ChatService_SubscribeChatroomServer) error {
	// Get or create the requested chatroom
	var chatroom chatdata.Chatroom
	{
		s.manager.GetLock().Lock()
		defer s.manager.GetLock().Unlock()
		chatroom = s.manager.GetOrCreateRoom(req.Chatroom.Name)
	}

	subscription := chatdata.NewSubscription(req.Self.Username)
	// Trigger this subscription to send an initial update
	subscription.SignalUpdate()

	// Add this subscription to the chatroom's current subscriptions
	{
		chatroom.GetLock().Lock()
		defer chatroom.GetLock().Unlock()
		chatroom.AddSubscription(subscription)
	}

	// Remove this subscription from the chatroom's current subscriptions at exit
	defer func() {
		chatroom.GetLock().Lock()
		defer chatroom.GetLock().Unlock()
		chatroom.RemoveSubscription(subscription.Id())
	}()

	for {
		select {
		case <-subscription.ShouldUpdate():
			if err := s.sendSubscriptionUpdate(stream); err != nil {
				return err
			}
		case <-stream.Context().Done():
			return nil
		}
	}
}

func Start(serverAddress string) error {
	lis, err := net.Listen("tcp", serverAddress)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		// Keepalive will disconnect an unresponsive client after approximately 1 minute (Time + Timeout).
		Time:    30 * time.Second,
		Timeout: 30 * time.Second,
	}))

	pb.RegisterChatServiceServer(grpcServer, newChatServer())

	fmt.Println("Running server...")
	grpcServer.Serve(lis)

	return nil
}
