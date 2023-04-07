package ephemeralstate

import (
	"log"
	"sync"

	pb "github.com/Richie78321/groupchat/chatservice"
	"google.golang.org/protobuf/proto"
)

const updatesBufferSize = 100

// ESGroup maps from server PID to its ephemeral state.
type ESGroup map[string]*pb.EphemeralState

type ESManager struct {
	Lock sync.RWMutex

	myPid   string
	esGroup ESGroup

	updates chan *pb.EphemeralState
}

func NewESManager(myPid string) *ESManager {
	return &ESManager{
		Lock:    sync.RWMutex{},
		myPid:   myPid,
		esGroup: make(map[string]*pb.EphemeralState),
		updates: make(chan *pb.EphemeralState, updatesBufferSize),
	}
}

func (e *ESManager) Updates() <-chan *pb.EphemeralState {
	return e.updates
}

// UpdateESLocked is the same as UpdateES except that it first acquires
// the ESManager lock.
func (e *ESManager) UpdateESLocked(pid string, es *pb.EphemeralState) {
	e.Lock.Lock()
	defer e.Lock.Unlock()

	e.UpdateES(pid, es)
}

func (e *ESManager) mergeEs(pid string, newEs *pb.EphemeralState) {
	if _, ok := e.esGroup[pid]; !ok {
		// There is no current ephemeral state for this PID, so no merging is
		// necessary.
		e.esGroup[pid] = newEs
		return
	}

	proto.Merge(e.esGroup[pid], newEs)
}

// signalChatrooms signals the chatrooms defined in this ephemeral state.
func (e *ESManager) signalChatrooms(es *pb.EphemeralState) {
	// TODO(richie): Implement
}

// UpdateES updates the ephemeral state for the server with the provided PID
// and triggers the related subscribers to update.
func (e *ESManager) UpdateES(pid string, newEs *pb.EphemeralState) {
	log.Printf("New update %v", newEs)

	if newEs == nil {
		panic("updated with nil ephemeral state")
	}

	e.mergeEs(pid, newEs)

	// Signal the chatrooms whose ephemeral state is being overwritten by the new
	// ephemeral state.
	e.signalChatrooms(newEs)

	if pid == e.myPid {
		// If this server's ephemeral state is updated, additionally broadcast the
		// ephemeral state update to peers.
		e.updates <- newEs
	}
}

// DeleteES removes the ephemeral state for the server with the provided PID
// and triggers the related subscribers to update.
// This is typically used when the connection to a peer is lost, so the ephemeral
// state from that peer is deleted.
func (e *ESManager) DeleteES(pid string) {
	if _, ok := e.esGroup[pid]; !ok {
		// The ephemeral state for this PID is already deleted, so there is nothing
		// to do.
		return
	}

	// Signal the chatrooms whose ephemeral state is being deleted.
	defer e.signalChatrooms(e.esGroup[pid])
	delete(e.esGroup, pid)
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
