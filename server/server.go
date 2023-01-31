package server

import (
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/Richie78321/groupchat/server/chatdata"
	"github.com/Richie78321/groupchat/server/chatdata/memory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type chatServer struct {
	manager chatdata.Manager
	pb.UnimplementedChatServiceServer
}

func newChatServer() *chatServer {
	return &chatServer{
		manager: memory.NewMemoryManager(),
	}
}

func sendSubscriptionUpdate(chatroom chatdata.Chatroom, stream pb.ChatService_SubscribeChatroomServer) error {
	chatroom.GetLock().Lock()
	defer chatroom.GetLock().Unlock()

	return stream.Send(&pb.ChatroomSubscriptionUpdate{
		Participants: chatroom.GetUsers(),
	})
}

func (s *chatServer) SubscribeChatroom(req *pb.SubscribeChatroomRequest, stream pb.ChatService_SubscribeChatroomServer) error {
	// Get or create the requested chatroom
	s.manager.GetLock().Lock()
	chatroom := s.manager.GetOrCreateRoom(req.Chatroom.Name)
	s.manager.GetLock().Unlock()

	subscription := chatdata.NewSubscription(req.Self)

	// Trigger this subscription to send an initial update
	subscription.SignalUpdate()

	// Add this subscription to the chatroom's current subscriptions
	chatroom.GetLock().Lock()
	chatroom.AddSubscription(subscription)
	chatroom.GetLock().Unlock()

	// Remove this subscription from the chatroom's current subscriptions at exit
	defer func() {
		chatroom.GetLock().Lock()
		chatroom.RemoveSubscription(subscription.Id())
		chatroom.GetLock().Unlock()
	}()

	fmt.Println("Accepted new client")

	for {
		select {
		case <-subscription.ShouldUpdate():
			if err := sendSubscriptionUpdate(chatroom, stream); err != nil {
				log.Printf("%v", err)
				return err
			}
		case <-stream.Context().Done():
			fmt.Println("Client disconnected")
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
