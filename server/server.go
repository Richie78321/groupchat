package server

import (
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
		// Assuming a reliable server, we use in-memory data
		manager: memory.NewMemoryManager(),
	}
}

// Send an update of the chatroom state over the provided stream.
func sendSubscriptionUpdate(chatroom chatdata.Chatroom, stream pb.ChatService_SubscribeChatroomServer) error {
	chatroom.GetLock().Lock()
	defer chatroom.GetLock().Unlock()

	return stream.Send(&pb.ChatroomSubscriptionUpdate{
		Participants: chatroom.Users(),
	})
}

func (s *chatServer) SubscribeChatroom(req *pb.SubscribeChatroomRequest, stream pb.ChatService_SubscribeChatroomServer) error {
	// Get or create the requested chatroom
	s.manager.GetLock().Lock()
	chatroom, ok := s.manager.Room(req.Chatroom.Name)
	if !ok {
		chatroom = s.manager.CreateRoom(req.Chatroom.Name)
	}
	s.manager.GetLock().Unlock()

	subscription := chatdata.NewSubscription(req.Self)

	// Trigger this subscription to send an initial update
	subscription.SignalUpdate()

	// Add this subscription to the chatroom's current subscriptions
	chatroom.GetLock().Lock()
	chatroom.AddSubscription(subscription)
	chatroom.GetLock().Unlock()
	log.Printf("Added subscription: user=%s, uuid=%v\n", req.Self.Username, subscription.Id())

	// Remove this subscription from the chatroom's current subscriptions at exit
	defer func() {
		chatroom.GetLock().Lock()
		chatroom.RemoveSubscription(subscription.Id())
		chatroom.GetLock().Unlock()
		log.Printf("Removed subscription: user=%s, uuid=%v\n", req.Self.Username, subscription.Id())
	}()

	for {
		select {
		case <-subscription.ShouldUpdate():
			// When signalled to update, send an update over the server stream
			if err := sendSubscriptionUpdate(chatroom, stream); err != nil {
				log.Printf("%v", err)
				return err
			}
		case <-stream.Context().Done():
			return nil
		}
	}
}

func Start(address string) error {
	// We strictly use TCP as the transport for reliable, in-order transfer.
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		// Keepalive will disconnect an unresponsive client after approximately 1 minute (Time + Timeout).
		// This means we have a maximum user online status staleness of around 1 minute.
		Time:    30 * time.Second,
		Timeout: 30 * time.Second,
	}))

	pb.RegisterChatServiceServer(grpcServer, newChatServer())

	log.Printf("Running server on %s...\n", address)
	grpcServer.Serve(lis)

	return nil
}
