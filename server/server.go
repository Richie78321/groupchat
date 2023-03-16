package server

import (
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/Richie78321/groupchat/chatservice"
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

	chatdata, err := sqlite.NewSqliteChatdata(fmt.Sprintf("%s.sqlite", id), id)
	if err != nil {
		return err
	}

	peerManager := replicationclient.NewPeerManager(peers, chatdata.NewEvents())
	replicationServer := replicationserver.NewReplicationServer(chatdata.EventsToBroadcast())
	// TODO(richie): Need to integrate SqliteChatdata with chatServer
	chatServer := chatserver.NewChatServer(peerManager)

	grpcServer := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		// Keepalive will disconnect an unresponsive client after approximately 1 minute (Time + Timeout).
		// This means we have a maximum user online status staleness of around 1 minute.
		Time:    30 * time.Second,
		Timeout: 30 * time.Second,
	}))
	pb.RegisterChatServiceServer(grpcServer, chatServer)
	pb.RegisterReplicationServiceServer(grpcServer, replicationServer)

	log.Printf("Running server on %s...\n", address)
	grpcServer.Serve(lis)

	return nil
}
