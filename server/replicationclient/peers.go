package replicationclient

import (
	"context"
	"io"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	pb "github.com/Richie78321/groupchat/chatservice"
)

type Peer struct {
	Id        string
	Addr      string
	Connected atomic.Bool
}

func NewPeer(id string, addr string) *Peer {
	return &Peer{
		Id:        id,
		Addr:      addr,
		Connected: atomic.Bool{},
	}
}

func (p *Peer) connect(m *PeerManager) {
	for {
		// Reset the state of the peer when retrying the connection
		m.esManager.DeleteESLocked(p.Id)
		p.Connected.Store(false)

		stream, err := p.attemptSubscribe(m)
		if err != nil {
			m.log.Printf("Failed to subscribe to `%s`: %v", p.Id, err)
			continue
		}

		// The peer is considered connected after successfully establishing an update subscription.
		m.log.Printf("Peer subscription to `%s` succeeeded", p.Id)
		p.Connected.Store(true)

		err = p.readUpdates(stream, m)
		if err != nil {
			m.log.Printf("Failed to read updates from subscription to `%s`: %v", p.Id, err)
			continue
		}

		m.log.Printf("Subscription to `%s` ended", p.Id)
	}
}

func (p *Peer) attemptSubscribe(m *PeerManager) (pb.ReplicationService_SubscribeUpdatesClient, error) {
	conn, err := grpc.Dial(
		p.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.WithBlock() is used to avoid returning from the handler before the connection has been
		// fully established.
		grpc.WithBlock(),
		// Keepalive will disconnect an unresponsive server after approximately 10 seconds (Time + Timeout).
		// This means we have a maximum peer connected status staleness of around 10 seconds.
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    5 * time.Second,
			Timeout: 5 * time.Second,
		}),
	)
	if err != nil {
		return nil, err
	}

	// Use the current sequence number vector. It is okay if this vector becomes
	// partially out-of-date due to new events, as duplicate incoming events are ignored.
	vector, err := m.synchronizer.SequenceNumberVector()
	if err != nil {
		return nil, err
	}

	client := pb.NewReplicationServiceClient(conn)
	stream, err := client.SubscribeUpdates(context.Background(), &pb.SubscribeRequest{
		SequenceNumberVector: vector,
	})
	if err != nil {
		return nil, err
	}

	return stream, nil
}

func (p *Peer) readUpdates(stream pb.ReplicationService_SubscribeUpdatesClient, m *PeerManager) error {
	for {
		update, err := stream.Recv()
		if err == io.EOF || stream.Context().Err() == context.Canceled {
			return nil
		}
		if err != nil {
			return err
		}

		m.log.Printf("Received update from `%s`", p.Id)

		if update.EphemeralState != nil {
			// A new ephemeral state has been received.
			// Update this peer's ephemeral state with the ESManager.
			m.esManager.UpdateESLocked(p.Id, update.EphemeralState)
		}

		// Update the garbage collection vector.
		if err := m.synchronizer.UpdateGarbageCollectedTo(update.GarbageCollectedToVector); err != nil {
			return err
		}

		for _, event := range update.Events {
			m.synchronizer.IncomingEvents() <- event
		}
	}
}
