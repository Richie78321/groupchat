package replicationclient

import "github.com/Richie78321/groupchat/server/chatdata"

type PeerManager struct {
	Peers        []*Peer
	synchronizer chatdata.EventSynchronizer
}

func NewPeerManager(peers []*Peer, synchronizer chatdata.EventSynchronizer) *PeerManager {
	p := &PeerManager{
		Peers:        peers,
		synchronizer: synchronizer,
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
