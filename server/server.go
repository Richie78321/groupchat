package server

import (
	"fmt"
	"net"
	"time"

	pb "github.com/Richie78321/groupchat/chatservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type chatServer struct {
	pb.UnimplementedChatServiceServer
}

func Start(serverAddress string) error {
	lis, err := net.Listen("tcp", serverAddress)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		// Keepalive will disconnect an unresponsive client after approximately 1 minute (Time + Timeout).
		Time:    30 * time.Second,
		Timeout: 30 * time.Second,
	}))

	pb.RegisterChatServiceServer(grpcServer, new(chatServer))

	fmt.Println("Running server...")
	grpcServer.Serve(lis)

	return nil
}
