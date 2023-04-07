package ephemeralstate

import (
	pb "github.com/Richie78321/groupchat/chatservice"
	"google.golang.org/protobuf/proto"
)

func (e *ESManager) ClientConnected(connected bool, chatroom string, user *pb.User) {
	e.Lock.Lock()
	defer e.Lock.Unlock()

	// Create a base chatroom ES that contains an empty ConnectedClients map.
	// This base proto is then merged with the existing chatroom ES (if it
	// exists), which clones the existing chatroom ES data.
	var baseChatroomEs *pb.ChatroomES = &pb.ChatroomES{
		ConnectedClients: make(map[string]*pb.User),
	}
	if es, ok := e.esGroup[e.myPid]; ok {
		if chatroomEs, ok := es.ChatroomEs[chatroom]; ok {
			proto.Merge(baseChatroomEs, chatroomEs)
		}
	}

	if _, currentlyConnected := baseChatroomEs.ConnectedClients[user.Username]; currentlyConnected == connected {
		// Client is already in the correct connection state,
		// so there is no need to update the ephemeral state.
		return
	}

	if connected {
		// Add the client to the set of connected clients.
		baseChatroomEs.ConnectedClients[user.Username] = user
	} else {
		// Remove the client from the set of connected clients.
		delete(baseChatroomEs.ConnectedClients, user.Username)
	}

	e.UpdateES(e.myPid, &pb.EphemeralState{
		ChatroomEs: map[string]*pb.ChatroomES{
			chatroom: baseChatroomEs,
		},
	})
}
