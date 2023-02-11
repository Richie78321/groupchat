package server

import (
	"log"
	"net"
	"time"

	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/Richie78321/groupchat/server/chatdata"
	"github.com/Richie78321/groupchat/server/chatdata/memory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
)

type chatServer struct {
	manager chatdata.Manager
	pb.UnimplementedChatServiceServer
}

func newChatServer() *chatServer {
	return &chatServer{
		// Assuming a reliable server, we use in-memory data structures
		manager: memory.NewMemoryManager(),
	}
}

// Return a gRPC error if the specified chatroom does not exist.
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
