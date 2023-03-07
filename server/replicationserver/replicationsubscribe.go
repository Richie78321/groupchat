package replicationserver

import (
	"log"

	pb "github.com/Richie78321/groupchat/chatservice"
)

const eventBufferSize = 100

type subscription struct {
	// A channel of events to broadcast.
	eventsToBroadcast chan []*pb.Event
}

func newSubscription() *subscription {
	return &subscription{
		eventsToBroadcast: make(chan []*pb.Event, eventBufferSize),
	}
}

func (s *subscription) broadcastEvents(events []*pb.Event) {
	s.eventsToBroadcast <- events
}

func sendSubscriptionUpdate(events []*pb.Event, stream pb.ReplicationService_SubscribeUpdatesServer) error {
	return stream.Send(&pb.SubscriptionUpdate{
		EphemeralState: &pb.EphemeralState{},
		Events:         events,
	})
}

func (s *ReplicationServer) SubscribeUpdates(req *pb.SubscribeRequest, stream pb.ReplicationService_SubscribeUpdatesServer) error {
	subscription := newSubscription()

	// Add this subscription to the current subscriptions
	s.lock.Lock()
	s.addSubscription(subscription)
	s.lock.Unlock()

	// Remove this subscription from the current subscriptions at exit
	defer func() {
		s.lock.Lock()
		s.removeSubscription(subscription)
		s.lock.Unlock()
	}()

	// TODO(richie): Need to trigger a special initial update that diffs events between the processes. The subscribe request should send a process event sequence number vector.

	log.Printf("Peer subscribed")

	for {
		select {
		case events := <-subscription.eventsToBroadcast:
			// When new events are available, send a subscription update
			if err := sendSubscriptionUpdate(events, stream); err != nil {
				log.Printf("%v", err)
				return err
			}
		case <-stream.Context().Done():
			return nil
		}
	}
}
