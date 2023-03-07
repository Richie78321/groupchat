package chatserver

import (
	"github.com/Richie78321/groupchat/server/chatdata"
	"github.com/Richie78321/groupchat/server/chatdata/memory"
	"github.com/Richie78321/groupchat/server/replication"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/Richie78321/groupchat/chatservice"
)

type ChatServer struct {
	manager     chatdata.Manager
	peerManager *replication.PeerManager
	pb.UnimplementedChatServiceServer
}

func NewChatServer(peerManager *replication.PeerManager) *ChatServer {
	return &ChatServer{
		// Assuming a reliable server, we use in-memory data structures
		manager:     memory.NewMemoryManager(),
		peerManager: peerManager,
	}
}

// Return a gRPC error if the specified chatroom does not exist.
func (s *ChatServer) getChatroomOrFail(r *pb.Chatroom) (chatdata.Chatroom, error) {
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
