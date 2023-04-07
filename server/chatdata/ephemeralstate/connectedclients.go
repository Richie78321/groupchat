package ephemeralstate

import (
	pb "github.com/Richie78321/groupchat/chatservice"
)

func (e *ESManager) ClientConnected(connected bool, chatroom string, user *pb.User) {
	e.Lock.Lock()
	defer e.Lock.Unlock()

	es, ok := e.ES()[e.myPid]
	if !ok {
		es = &pb.EphemeralState{
			ChatroomEs: make(map[string]*pb.ChatroomES),
		}
	}

	if _, ok := es.ChatroomEs[chatroom]; !ok {
		es.ChatroomEs[chatroom] = &pb.ChatroomES{
			ConnectedClients: make(map[string]*pb.User),
		}
	}
	chatroomEs := es.ChatroomEs[chatroom]

	if _, ok := chatroomEs.ConnectedClients[user.Username]; ok == connected {
		// Client is already in the correct connection state,
		// so there is no need to update the ephemeral state.
		return
	}

	if connected {
		// Add the client to the set of connected clients.
		chatroomEs.ConnectedClients[user.Username] = user
	} else {
		// Remove the client from the set of connected clients.
		delete(chatroomEs.ConnectedClients, user.Username)
	}

	e.UpdateES(e.myPid, es)
}
