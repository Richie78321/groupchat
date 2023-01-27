package server

import (
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/Richie78321/groupchat/chatservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type chatServer struct {
	pb.UnimplementedChatServiceServer
}

func (s *chatServer) TestStream(req *pb.TestRequest, stream pb.ChatService_TestStreamServer) error {
	fmt.Println("Waiting for client to be done")
	<-stream.Context().Done()

	fmt.Println("Client is done. Exiting")
	return nil
}

func Start() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:3000"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		Time: time.Minute,
	}))

	pb.RegisterChatServiceServer(grpcServer, new(chatServer))

	fmt.Println("Running server...")
	grpcServer.Serve(lis)
}
