package replicationserver

import (
	pb "github.com/Richie78321/groupchat/chatservice"
)

const (
	eventBufferSize = 100
	esBufferSize    = 100
)

type subscription struct {
	// A channel of events to broadcast.
	eventsToBroadcast chan *pb.Event

	// A channel of ephemeral state updates to broadcast.
	esToBroadcast chan *pb.EphemeralState
}

func newSubscription() *subscription {
	return &subscription{
		eventsToBroadcast: make(chan *pb.Event, eventBufferSize),
		esToBroadcast:     make(chan *pb.EphemeralState, esBufferSize),
	}
}

func (s *subscription) broadcastEvent(event *pb.Event) {
	s.eventsToBroadcast <- event
}

func (s *ReplicationServer) sendSubscriptionUpdate(events []*pb.Event, es *pb.EphemeralState, stream pb.ReplicationService_SubscribeUpdatesServer) error {
	return stream.Send(&pb.SubscriptionUpdate{
		EphemeralState:           es,
		Events:                   events,
		GarbageCollectedToVector: s.synchronizer.GarbageCollectedTo(),
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

	s.log.Printf("Peer subscribed with sequence number vector: %v", req.SequenceNumberVector)

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
	// In the initial subscription update, send the entirety of the current ephemeral state.
	// After this, only ephemeral state diffs will be sent.
	if err := s.sendSubscriptionUpdate(eventDiff, s.esManager.MyESLocked(), stream); err != nil {
		s.log.Printf("%v", err)
		return err
	}
	s.log.Print("Sent initial update to subscriber")

	for {
		// Batching multiple ephemeral state and event updates into a single subscription update
		// would be more efficient here. This can be revisited in the future if necessary.
		select {
		case event := <-subscription.eventsToBroadcast:
			// When new events are available, send a subscription update.
			if err := s.sendSubscriptionUpdate([]*pb.Event{event}, nil, stream); err != nil {
				s.log.Printf("%v", err)
				return err
			}
		case newEs := <-subscription.esToBroadcast:
			// When a new ephemeral state has been set, send a subscription update.
			if err := s.sendSubscriptionUpdate([]*pb.Event{}, newEs, stream); err != nil {
				s.log.Printf("%v", err)
				return err
			}
		case <-stream.Context().Done():
			return nil
		}
	}
}
