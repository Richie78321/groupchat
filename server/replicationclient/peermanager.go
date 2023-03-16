package replicationclient

import pb "github.com/Richie78321/groupchat/chatservice"

type PeerManager struct {
	Peers     []*Peer
	newEvents chan<- *pb.Event
}

func NewPeerManager(peers []*Peer, newEvents chan<- *pb.Event) *PeerManager {
	p := &PeerManager{
		Peers:     peers,
		newEvents: newEvents,
	}

	// Spawn threads to manage peer connections
	go p.spawnPeerThreads()

	return p
}

func (m *PeerManager) spawnPeerThreads() {
	// Spawn a thread to manage connections to each peer
	for _, peer := range m.Peers {
		// TODO(richie): Potentially use context here to make threads cancellable
		go peer.connect(m)
	}
}
