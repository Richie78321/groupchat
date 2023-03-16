package replicationclient

import (
	"context"
	"log"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	pb "github.com/Richie78321/groupchat/chatservice"
)

// EphemeralStateHolder is necessary to avoid storing a nil value in atomic.Value when
// *pb.EphemeralState can sometimes be nil.
type EphemeralStateHolder struct {
	state *pb.EphemeralState
}

type Peer struct {
	Id             string
	Addr           string
	Connected      atomic.Bool
	EphemeralState atomic.Value
}

func NewPeer(id string, addr string) *Peer {
	return &Peer{
		Id:             id,
		Addr:           addr,
		Connected:      atomic.Bool{},
		EphemeralState: atomic.Value{},
	}
}

func (p *Peer) connect(m *PeerManager) {
	for {
		// Reset the state of the peer when retrying the connection
		p.EphemeralState.Store(&EphemeralStateHolder{state: nil})
		p.Connected.Store(false)

		stream, err := p.attemptSubscribe(m)
		if err != nil {
			log.Printf("Failed to subscribe to `%s`: %v", p.Id, err)
			continue
		}

		// The peer is considered connected after successfully establishing an update subscription.
		log.Printf("Peer subscription to `%s` succeeeded", p.Id)
		p.Connected.Store(true)

		err = p.readUpdates(stream, m)
		if err != nil {
			log.Printf("Failed to read updates from subscription to `%s`: %v", p.Id, err)
			continue
		}
	}
}

func (p *Peer) attemptSubscribe(m *PeerManager) (pb.ReplicationService_SubscribeUpdatesClient, error) {
	conn, err := grpc.Dial(
		p.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.WithBlock() is used to avoid returning from the handler before the connection has been
		// fully established.
		grpc.WithBlock(),
		// Keepalive will disconnect an unresponsive server after approximately 1 minute (Time + Timeout).
		// This means we have a maximum peer connected status staleness of around 1 minute.
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    30 * time.Second,
			Timeout: 30 * time.Second,
		}),
	)
	if err != nil {
		return nil, err
	}

	// Retrieve the current sequence number vector. It is okay if this vector becomes
	// partially out-of-date due to new events, as duplicate events are ignored.

	client := pb.NewReplicationServiceClient(conn)
	stream, err := client.SubscribeUpdates(context.Background(), &pb.SubscribeRequest{})
	if err != nil {
		return nil, err
	}

	return stream, nil
}

func (p *Peer) readUpdates(stream pb.ReplicationService_SubscribeUpdatesClient, m *PeerManager) error {
	for {
		update, err := stream.Recv()
		if err != nil {
			return err
		}

		log.Printf("Received update from `%s`", p.Id)

		// Update ephemeral state before delivering the event.
		p.EphemeralState.Store(&EphemeralStateHolder{
			state: update.EphemeralState,
		})

		for _, event := range update.Events {
			m.synchronizer.IncomingEvents() <- event
		}
	}
}
