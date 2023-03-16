package replicationserver

import (
	"log"

	pb "github.com/Richie78321/groupchat/chatservice"
)

const eventBufferSize = 100

type subscription struct {
	// A channel of events to broadcast.
	eventsToBroadcast chan *pb.Event
}

func newSubscription() *subscription {
	return &subscription{
		eventsToBroadcast: make(chan *pb.Event, eventBufferSize),
	}
}

func (s *subscription) broadcastEvent(event *pb.Event) {
	s.eventsToBroadcast <- event
}

func sendSubscriptionUpdate(event *pb.Event, stream pb.ReplicationService_SubscribeUpdatesServer) error {
	return stream.Send(&pb.SubscriptionUpdate{
		// TODO(richie): Send a real ephemeral state here
		EphemeralState: &pb.EphemeralState{},
		Events:         []*pb.Event{event},
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

	log.Printf("Peer subscribed")

	// The initial update must happen strictly after the subscription is registered.
	// Otherwise events could be lost in the time between the initial update and subscription registration.

	// TODO(richie): Need to trigger a special initial update that diffs events between the processes. The subscribe request should send a process event sequence number vector.

	for {
		select {
		case event := <-subscription.eventsToBroadcast:
			// When new events are available, send a subscription update
			if err := sendSubscriptionUpdate(event, stream); err != nil {
				log.Printf("%v", err)
				return err
			}
		case <-stream.Context().Done():
			return nil
		}
	}
}
