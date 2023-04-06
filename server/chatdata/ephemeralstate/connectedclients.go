package ephemeralstate

import (
	pb "github.com/Richie78321/groupchat/chatservice"
	"google.golang.org/protobuf/proto"
)

func (e *ESManager) ClientConnected(connected bool, user *pb.User) {
	e.Lock.Lock()
	defer e.Lock.Unlock()

	var baseEs *pb.EphemeralState
	if myEs, esOk := e.ES()[e.myPid]; esOk {
		// The new ephemeral state should be based on a copy of
		// the current.
		baseEs = proto.Clone(myEs).(*pb.EphemeralState)
	} else {
		baseEs = &pb.EphemeralState{}
	}

	if baseEs.ConnectedClients == nil {
		baseEs.ConnectedClients = make(map[string]*pb.User)
	}

	if _, ok := baseEs.ConnectedClients[user.Username]; ok == connected {
		// Client is already in the correct connection state,
		// so there is no need to update the ephemeral state.
		return
	}

	if connected {
		// Add the client to the set of connected clients.
		baseEs.ConnectedClients[user.Username] = user
	} else {
		// Remove the client from the set of connected clients.
		delete(baseEs.ConnectedClients, user.Username)
	}

	e.UpdateES(e.myPid, baseEs)
}
