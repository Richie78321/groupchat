package replicationserver

import (
	"log"

	pb "github.com/Richie78321/groupchat/chatservice"
)

type subscription struct {
	// A channel for signalling updates. Needs a strict buffer size of 1.
	update chan struct{}
}

func newSubscription() *subscription {
	return &subscription{
		// The buffer size is strictly 1 to ensure proper signalling behavior.
		update: make(chan struct{}, 1),
	}
}

func (s *subscription) updateSubscription() {
	s.update <- struct{}{}
}

func sendSubscriptionUpdate(stream pb.ReplicationService_SubscribeUpdatesServer) error {
	return stream.Send(&pb.SubscriptionUpdate{
		EphemeralState: &pb.EphemeralState{},
		Events:         make([]*pb.Event, 0),
	})
}

func (s *ReplicationServer) SubscribeUpdates(req *pb.SubscribeRequest, stream pb.ReplicationService_SubscribeUpdatesServer) error {
	subscription := newSubscription()

	// Add this subscription to the current subscriptions
	s.lock.Lock()
	s.addSubscription(subscription)
	s.lock.Unlock()

	log.Printf("Peer subscribed")

	// Remove this subscription from the current subscriptions at exit
	defer func() {
		s.lock.Lock()
		s.removeSubscription(subscription)
		s.lock.Unlock()
	}()

	for {
		select {
		case <-subscription.update:
			// When there are new events, send an update over the server stream
			if err := sendSubscriptionUpdate(stream); err != nil {
				log.Printf("%v", err)
				return err
			}
		case <-stream.Context().Done():
			return nil
		}
	}
}
