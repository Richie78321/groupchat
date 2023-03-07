package replicationserver

import (
	"sync"

	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/Richie78321/groupchat/server/replicationclient"
)

type ReplicationServer struct {
	lock sync.Mutex

	subscriptions map[*subscription]struct{}

	peerManager *replicationclient.PeerManager
	pb.UnimplementedReplicationServiceServer
}

func (r *ReplicationServer) SignalSubscriptions() {
	r.lock.Lock()
	defer r.lock.Unlock()

	for subscription := range r.subscriptions {
		subscription.signalUpdate()
	}
}

func (r *ReplicationServer) addSubscription(s *subscription) {
	r.subscriptions[s] = struct{}{}
}

func (r *ReplicationServer) removeSubscription(s *subscription) {
	delete(r.subscriptions, s)
}

func NewReplicationServer(peerManager *replicationclient.PeerManager) *ReplicationServer {
	return &ReplicationServer{
		lock:          sync.Mutex{},
		subscriptions: make(map[*subscription]struct{}),
		peerManager:   peerManager,
	}
}
