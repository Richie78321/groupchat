package replicationserver

import (
	"sync"

	pb "github.com/Richie78321/groupchat/chatservice"
)

type ReplicationServer struct {
	lock sync.Mutex

	subscriptions     map[*subscription]struct{}
	eventsToBroadcast <-chan *pb.Event

	pb.UnimplementedReplicationServiceServer
}

func NewReplicationServer(eventsToBroadcast <-chan *pb.Event) *ReplicationServer {
	r := &ReplicationServer{
		lock:              sync.Mutex{},
		subscriptions:     make(map[*subscription]struct{}),
		eventsToBroadcast: eventsToBroadcast,
	}

	// Spawn a thread to broadcast events to subscriptions
	go r.broadcastEvents()

	return r
}

func (r *ReplicationServer) broadcastEvents() {
	for {
		event := <-r.eventsToBroadcast

		r.lock.Lock()
		for subscription := range r.subscriptions {
			subscription.broadcastEvent(event)
		}
		r.lock.Unlock()
	}
}

func (r *ReplicationServer) addSubscription(s *subscription) {
	r.subscriptions[s] = struct{}{}
}

func (r *ReplicationServer) removeSubscription(s *subscription) {
	delete(r.subscriptions, s)
}
