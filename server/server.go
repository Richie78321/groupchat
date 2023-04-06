package server

import (
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/Richie78321/groupchat/server/chatdata/ephemeralstate"
	"github.com/Richie78321/groupchat/server/chatdata/sqlite"
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

	peerPids := make([]string, len(peers))
	for i, peer := range peers {
		peerPids[i] = peer.Id
	}

	chatdata, err := sqlite.NewSqliteChatdata(fmt.Sprintf("./%s.sqlite", id), fmt.Sprintf("./%s.json", id), id, peerPids)
	if err != nil {
		return err
	}
	esManager := ephemeralstate.NewESManager(id)

	chatdataManager := sqlite.NewChatdataManager(chatdata, esManager)
	chatdata.SubscriptionSignal = chatdataManager

	peerManager := replicationclient.NewPeerManager(peers, esManager, chatdata)
	replicationServer := replicationserver.NewReplicationServer(chatdata)
	chatServer := chatserver.NewChatServer(chatdataManager, peerManager)

	grpcServer := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		// Keepalive will disconnect an unresponsive client after approximately 1 minute (Time + Timeout).
		// This means we have a maximum user online status staleness of around 1 minute.
		Time:    30 * time.Second,
		Timeout: 30 * time.Second,
	}), grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		// Permit keepalive pings from the client at most every 30 seconds. This allows
		// the client to send frequent keepalive pings.
		MinTime: 30 * time.Second,
	}))
	pb.RegisterChatServiceServer(grpcServer, chatServer)
	pb.RegisterReplicationServiceServer(grpcServer, replicationServer)

	log.Printf("Running server on %s...\n", address)
	grpcServer.Serve(lis)

	return nil
}
