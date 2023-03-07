package replication

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

const (
	connectionTimeout    time.Duration = 10 * time.Second
	connectionRetryDelay time.Duration = 10 * time.Second
)

type PeerManager struct {
	Peers  []*Peer
	Events chan struct{}
}

func NewPeerManager(peers []*Peer) *PeerManager {
	return &PeerManager{
		Peers:  peers,
		Events: make(chan struct{}),
	}
}

func (m *PeerManager) ConnectPeers() {
	// Spawn a thread to manage connections to each peer
	for _, peer := range m.Peers {
		// TODO(richie): Potentially use context here to make threads cancellable
		go peer.connect(m.Events)
	}
}

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

func (p *Peer) connect(events chan<- struct{}) {
	for {
		p.Connected.Store(false)

		stream, err := p.attemptSubscribe()
		if err != nil {
			log.Printf("Failed to subscribe to `%s`: %v", p.Id, err)
			continue
		}

		// The peer is considered connected after successfully establishing an update subscription.
		log.Printf("Peer subscription to `%s` succeeeded", p.Id)
		p.Connected.Store(true)

		err = p.readUpdates(stream, events)
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

func (p *Peer) readUpdates(stream pb.ReplicationService_SubscribeUpdatesClient, events chan<- struct{}) error {
	for {
		_, err := stream.Recv()
		if err != nil {
			return err
		}

		// TODO(richie): Replace with real updates when they become available.
		events <- struct{}{}
	}
}
