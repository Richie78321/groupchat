package server

import (
	"log"
	"time"

	pb "github.com/Richie78321/groupchat/chatservice"
)

func (s *chatServer) SubscribeUpdates(req *pb.SubscribeRequest, stream pb.ReplicationService_SubscribeUpdatesServer) error {
	// TODO(richie): Update this to something meaningful
	log.Printf("Peer subscribed")
	time.Sleep(10 * time.Hour)
	return nil
}
