package server

import (
	"fmt"
	"log"
	"net"

	pb "github.com/Richie78321/groupchat/chatservice"
	"google.golang.org/grpc"
)

type chatServer struct {
	pb.UnimplementedChatServiceServer
}

func (s *chatServer) TestStream(req *pb.TestRequest, stream pb.ChatService_TestStreamServer) error {
	return nil
}

func Start() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:3000"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterChatServiceServer(grpcServer, new(chatServer))

	fmt.Println("Running server...")
	grpcServer.Serve(lis)
}
