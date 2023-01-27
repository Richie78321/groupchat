package client

import (
	"context"
	"log"

	pb "github.com/Richie78321/groupchat/chatservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Start() {
	conn, err := grpc.Dial("localhost:3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	client := pb.NewChatServiceClient(conn)
	stream, err := client.TestStream(context.Background(), &pb.TestRequest{})
	if err != nil {
		log.Fatalf("failed to send request: %v", err)
	}

	<-stream.Context().Done()
}
