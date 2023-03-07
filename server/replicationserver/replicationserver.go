package replicationserver

import (
	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/Richie78321/groupchat/server/replication"
)

type ReplicationServer struct {
	peerManager *replication.PeerManager
	pb.UnimplementedReplicationServiceServer
}

func NewReplicationServer(peerManager *replication.PeerManager) *ReplicationServer {
	return &ReplicationServer{
		peerManager: peerManager,
	}
}
