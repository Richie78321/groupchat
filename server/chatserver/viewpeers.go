package chatserver

import (
	"context"

	pb "github.com/Richie78321/groupchat/chatservice"
)

func (s *ChatServer) ViewPeers(ctx context.Context, req *pb.ViewPeersRequest) (*pb.ViewPeersResponse, error) {
	// Only include the peers that the server is currently connected to.
	connectedPeers := make([]*pb.Peer, 0)
	for _, peer := range s.peerManager.Peers {
		if !peer.Connected.Load() {
			continue
		}

		connectedPeers = append(connectedPeers, &pb.Peer{
			Id: peer.Id,
		})
	}

	return &pb.ViewPeersResponse{
		Peers: connectedPeers,
	}, nil
}
