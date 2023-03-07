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

type PeerManager struct {
	Peers            []*Peer
	delivered_events chan *pb.Event
}

func NewPeerManager(peers []*Peer) *PeerManager {
	return &PeerManager{
		Peers:            peers,
		delivered_events: make(chan *pb.Event),
	}
}

func (m *PeerManager) deliverEvent(e *pb.Event) {
	m.delivered_events <- e
}

func (m *PeerManager) ConnectPeers() {
	// Spawn a thread to manage connections to each peer
	for _, peer := range m.Peers {
		// TODO(richie): Potentially use context here to make threads cancellable
		go peer.connect(m)
	}
}

func (m *PeerManager) Events() <-chan *pb.Event {
	return m.delivered_events
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
		p.EphemeralState.Store(nil)
		p.Connected.Store(false)

		stream, err := p.attemptSubscribe()
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

func (p *Peer) attemptSubscribe() (pb.ReplicationService_SubscribeUpdatesClient, error) {
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
		p.EphemeralState.Store(update.EphemeralState)

		for _, event := range update.Events {
			m.deliverEvent(event)
		}
	}
}
