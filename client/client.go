package client

import (
	"context"
	"fmt"
	"io"

	pb "github.com/Richie78321/groupchat/chatservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Start(target string) error {
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	client := pb.NewChatServiceClient(conn)
	stream, err := client.SubscribeChatroom(context.TODO(), &pb.SubscribeChatroomRequest{
		Self: &pb.User{
			Username: "richie",
		},
		Chatroom: &pb.Chatroom{
			Name: "test",
		},
	})
	if err != nil {
		return err
	}

	num := 0
	for {
		fmt.Println("Waiting for next update...")
		update, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		num += 1
		fmt.Printf("Update %d: %s\n", num, update.String())
	}

	return nil
}
