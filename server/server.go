package server

import (
	"log"
	"net"
	"time"

	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/Richie78321/groupchat/server/chatdata"
	"github.com/Richie78321/groupchat/server/chatdata/memory"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
)

const latestMessageWindow = 10

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

	latestMessagesPb := make([]*pb.Message, latestMessageWindow)
	for i, m := range chatroom.GetLatestMessages(latestMessageWindow) {
		latestMessagesPb[i] = chatdata.MessageToPb(m)
	}

	return stream.Send(&pb.ChatroomSubscriptionUpdate{
		Participants:   chatroom.Users(),
		LatestMessages: latestMessagesPb,
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

func (s *chatServer) getChatroomOrFail(r *pb.Chatroom) (chatdata.Chatroom, error) {
	chatroom, ok := s.manager.Room(r.Name)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "chatroom `%s` does not exist", r.Name)
	}

	return chatroom, nil
}

// Return a gRPC error if the user is not logged into the chatroom.
func ensureUserLoggedIn(c chatdata.Chatroom, u *pb.User) error {
	// There is currently no way to ensure that the request is actually being made
	// by this user. This could be improved by using subscription UUIDs (which are
	// hard to guess) as the identifier. We avoid this for simplicity.
	users := c.Users()

	// This could be made faster than a linear scan, but this works for now.
	for _, user := range users {
		if user.Username == u.Username {
			return nil
		}
	}

	return status.Errorf(codes.PermissionDenied, "user `%s` is not logged into chatroom `%s`", u.Username, c.RoomName())
}

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
