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

func (s *subscription) signalUpdate() {
	select {
	case s.update <- struct{}{}:
	default:
		// The update channel has a strict buffer size of 1.
		// If there is already a signal in the subscription update channel,
		// then we can safely continue because the subscription has already
		// been signalled to update.
	}
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
		case <-subscription.update:
			// When an update signal is made, send an update over the server stream
			if err := sendSubscriptionUpdate(stream); err != nil {
				log.Printf("%v", err)
				return err
			}
		case <-stream.Context().Done():
			return nil
		}
	}
}
