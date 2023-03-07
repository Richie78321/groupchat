package replicationserver

import (
	"sync"

	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/Richie78321/groupchat/server/replication"
)

type ReplicationServer struct {
	lock sync.Mutex

	subscriptions map[*subscription]struct{}

	peerManager *replication.PeerManager
	pb.UnimplementedReplicationServiceServer
}

func (r *ReplicationServer) BroadcastEvents(events []*pb.Event) {
	r.lock.Lock()
	defer r.lock.Unlock()

	for subscription := range r.subscriptions {
		subscription.broadcastEvents(events)
	}
}

func (r *ReplicationServer) addSubscription(s *subscription) {
	r.subscriptions[s] = struct{}{}
}

func (r *ReplicationServer) removeSubscription(s *subscription) {
	delete(r.subscriptions, s)
}

func NewReplicationServer(peerManager *replication.PeerManager) *ReplicationServer {
	return &ReplicationServer{
		lock:          sync.Mutex{},
		subscriptions: make(map[*subscription]struct{}),
		peerManager:   peerManager,
	}
}
