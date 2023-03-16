package server

import (
	"log"
	"net"
	"time"

	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/Richie78321/groupchat/server/chatserver"
	"github.com/Richie78321/groupchat/server/replicationclient"
	"github.com/Richie78321/groupchat/server/replicationserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func Start(id string, address string, peers []*replicationclient.Peer) error {
	// We strictly use TCP as the transport for reliable, in-order transfer.
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		// Keepalive will disconnect an unresponsive client after approximately 1 minute (Time + Timeout).
		// This means we have a maximum user online status staleness of around 1 minute.
		Time:    30 * time.Second,
		Timeout: 30 * time.Second,
	}))

	peerManager := replicationclient.NewPeerManager(peers)
	chatServer := chatserver.NewChatServer(peerManager)
	replicationServer := replicationserver.NewReplicationServer(peerManager)

	pb.RegisterChatServiceServer(grpcServer, chatServer)
	pb.RegisterReplicationServiceServer(grpcServer, replicationServer)

	log.Printf("Running server on %s...\n", address)
	peerManager.SpawnPeerThreads()
	grpcServer.Serve(lis)

	return nil
}
