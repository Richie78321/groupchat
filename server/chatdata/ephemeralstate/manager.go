package ephemeralstate

import (
	"sync"

	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/Richie78321/groupchat/server/util"
)

// ESGroup maps from server PID to its ephemeral state.
type ESGroup map[string]*pb.EphemeralState

type ESManager struct {
	Lock sync.RWMutex

	myPid   string
	esGroup ESGroup

	Update *util.Signal
}

func NewESManager(myPid string) *ESManager {
	return &ESManager{
		Lock:    sync.RWMutex{},
		myPid:   myPid,
		esGroup: make(map[string]*pb.EphemeralState),
		Update:  util.NewSignal(),
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
		e.Update.Signal()
	}
}

func (e *ESManager) ES() ESGroup {
	return e.esGroup
}

// MyESLocked is the same as MyES except that it first acquires
// the ESManager lock.
func (e *ESManager) MyESLocked() *pb.EphemeralState {
	e.Lock.RLock()
	defer e.Lock.RUnlock()

	return e.MyES()
}

func (e *ESManager) MyES() *pb.EphemeralState {
	return e.esGroup[e.myPid]
}
