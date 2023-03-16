package replicationserver

import (
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

func sendSubscriptionUpdate(events []*pb.Event, stream pb.ReplicationService_SubscribeUpdatesServer) error {
	return stream.Send(&pb.SubscriptionUpdate{
		// TODO(richie): Send a real ephemeral state here
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

	s.log.Printf("Peer subscribed")

	// Send an initial subscription update that includes the events the subscriber does not
	// already have, according to the attached sequence number vector.
	//
	// This initial update must happen strictly after the subscription is registered.
	// Otherwise events could be lost in the time between the initial update and subscription registration.
	eventDiff, err := s.synchronizer.EventDiff(req.SequenceNumberVector)
	if err != nil {
		s.log.Printf("%v", err)
		return err
	}
	if err := sendSubscriptionUpdate(eventDiff, stream); err != nil {
		s.log.Printf("%v", err)
		return err
	}

	for {
		select {
		case event := <-subscription.eventsToBroadcast:
			// When new events are available, send a subscription update
			if err := sendSubscriptionUpdate([]*pb.Event{event}, stream); err != nil {
				s.log.Printf("%v", err)
				return err
			}
		case <-stream.Context().Done():
			return nil
		}
	}
}
