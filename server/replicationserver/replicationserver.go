package replicationserver

import (
	"log"
	"os"
	"sync"

	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/Richie78321/groupchat/server/chatdata"
	"github.com/Richie78321/groupchat/server/chatdata/ephemeralstate"
)

type ReplicationServer struct {
	lock sync.Mutex

	subscriptions map[*subscription]struct{}
	synchronizer  chatdata.EventSynchronizer

	esManager *ephemeralstate.ESManager

	log *log.Logger

	pb.UnimplementedReplicationServiceServer
}

func NewReplicationServer(synchronizer chatdata.EventSynchronizer, esManager *ephemeralstate.ESManager) *ReplicationServer {
	r := &ReplicationServer{
		lock:          sync.Mutex{},
		subscriptions: make(map[*subscription]struct{}),
		synchronizer:  synchronizer,
		esManager:     esManager,
		log:           log.New(os.Stdout, "[Replication Server] ", log.Default().Flags()),
	}

	// Spawn a thread to broadcast events to subscriptions
	go r.broadcastEvents()

	// Spawn a thread to broadcast ephemeral state update signals to subscriptions
	go r.signalESUpdates()

	return r
}

func (r *ReplicationServer) broadcastEvents() {
	for {
		event := <-r.synchronizer.OutgoingEvents()
		r.log.Printf("Broadcasting event PID=%s, SEQ=%d, LTS=%d", event.Pid, event.SequenceNumber, event.LamportTimestamp)

		r.lock.Lock()
		for subscription := range r.subscriptions {
			subscription.broadcastEvent(event)
		}
		r.lock.Unlock()
	}
}

func (r *ReplicationServer) signalESUpdates() {
	for {
		newEs := <-r.esManager.Updates()
		r.log.Print("Broadcasting ES update")

		r.lock.Lock()
		for subscription := range r.subscriptions {
			subscription.esToBroadcast <- newEs
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
