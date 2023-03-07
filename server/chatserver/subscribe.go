package chatserver

import (
	"log"

	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/Richie78321/groupchat/server/chatdata"
)

const latestMessageWindow = 10

// Send an update of the chatroom state over the provided stream.
func sendSubscriptionUpdate(chatroom chatdata.Chatroom, stream pb.ChatService_SubscribeChatroomServer) error {
	chatroom.GetLock().Lock()
	defer chatroom.GetLock().Unlock()

	return stream.Send(&pb.ChatroomSubscriptionUpdate{
		Participants:   chatroom.Users(),
		LatestMessages: chatdata.MessageListToPb(chatroom.LatestMessages(latestMessageWindow)),
	})
}

func (s *ChatServer) SubscribeChatroom(req *pb.SubscribeChatroomRequest, stream pb.ChatService_SubscribeChatroomServer) error {
	// Get or create the requested chatroom
	s.manager.GetLock().Lock()
	chatroom, ok := s.manager.Room(req.Chatroom.Name)
	if !ok {
		chatroom = s.manager.CreateRoom(req.Chatroom.Name)
	}
	s.manager.GetLock().Unlock()

	subscription := chatdata.NewSubscription(req.Self)

	// Trigger this subscription to send an initial update
	subscription.SignalUpdate()

	chatroom.GetLock().Lock()
	// Add this subscription to the chatroom's current subscriptions
	chatroom.AddSubscription(subscription)
	// Signal the subscriptions because a new user has logged into the chatroom
	chatroom.SignalSubscriptions()
	chatroom.GetLock().Unlock()
	log.Printf("Added subscription: user=%s, uuid=%v\n", req.Self.Username, subscription.Id())

	// Remove this subscription from the chatroom's current subscriptions at exit
	defer func() {
		chatroom.GetLock().Lock()
		chatroom.RemoveSubscription(subscription.Id())
		// Signal the subscriptions because a user has logged out of the chatroom
		chatroom.SignalSubscriptions()
		chatroom.GetLock().Unlock()
		log.Printf("Removed subscription: user=%s, uuid=%v\n", req.Self.Username, subscription.Id())
	}()

	for {
		select {
		case <-subscription.ShouldUpdate():
			// When signalled to update, send an update over the server stream
			if err := sendSubscriptionUpdate(chatroom, stream); err != nil {
				log.Printf("%v", err)
				return err
			}
		case <-stream.Context().Done():
			return nil
		}
	}
}
