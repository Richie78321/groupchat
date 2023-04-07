package replicationclient

import (
	"log"
	"os"

	"github.com/Richie78321/groupchat/server/chatdata"
	"github.com/Richie78321/groupchat/server/chatdata/ephemeralstate"
)

type PeerManager struct {
	Peers        []*Peer
	synchronizer chatdata.EventSynchronizer

	esManager *ephemeralstate.ESManager

	log *log.Logger
}

func NewPeerManager(peers []*Peer, esManager *ephemeralstate.ESManager, synchronizer chatdata.EventSynchronizer) *PeerManager {
	p := &PeerManager{
		Peers:        peers,
		synchronizer: synchronizer,

		esManager: esManager,

		log: log.New(os.Stdout, "[Replication Client] ", log.Default().Flags()),
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
