package ephemeralstate

import (
	"sync"

	pb "github.com/Richie78321/groupchat/chatservice"
)

// ESGroup maps from server PID to its ephemeral state.
type ESGroup map[string]*pb.EphemeralState

type ESManager struct {
	Lock sync.RWMutex

	myPid   string
	esGroup ESGroup
}

func NewESManager(myPid string) *ESManager {
	return &ESManager{
		Lock:    sync.RWMutex{},
		myPid:   myPid,
		esGroup: make(map[string]*pb.EphemeralState),
	}
}

// UpdateESLocked is the same as UpdateES except that it first acquires
// the ESManager lock.
func (e *ESManager) UpdateESLocked(pid string, es *pb.EphemeralState) {
	e.Lock.Lock()
	defer e.Lock.Unlock()

	e.UpdateES(pid, es)
}

func (e *ESManager) UpdateES(pid string, es *pb.EphemeralState) {
	e.esGroup[pid] = es
	if pid == e.myPid {
		// TODO(richie): Notify peers when your ES updates.
	}
}

func (e *ESManager) ES() ESGroup {
	return e.esGroup
}
