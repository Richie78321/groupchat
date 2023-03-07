package replicationserver

import (
	"log"

	pb "github.com/Richie78321/groupchat/chatservice"
)

type subscription struct {
	events chan []*pb.Event
}

func newSubscription() *subscription {
	return &subscription{
		events: make(chan []*pb.Event),
	}
}

func (s *subscription) broadcastEvents(events []*pb.Event) {
	// Place the events on the events channel to trigger a subscription update
	s.events <- events
}

func sendSubscriptionUpdate(newEvents []*pb.Event, stream pb.ReplicationService_SubscribeUpdatesServer) error {
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
		case newEvents := <-subscription.events:
			// When there are new events, send an update over the server stream
			if err := sendSubscriptionUpdate(newEvents, stream); err != nil {
				log.Printf("%v", err)
				return err
			}
		case <-stream.Context().Done():
			return nil
		}
	}
}
