package client

import (
	"context"
	"fmt"

	pb "github.com/Richie78321/groupchat/chatservice"
)

type sendArgs struct {
	Args struct {
		Message string `description:"message to send"`
	} `positional-args:"yes" required:"yes"`
}

func init() {
	parser.AddCommand("a", "send a message", "", &sendArgs{})
}

func (s *sendArgs) Execute(args []string) error {
	if len(s.Args.Message) <= 0 {
		fmt.Println("Message cannot be empty")
		return nil
	}

	if !inChatroom() {
		fmt.Println("Not in a chatroom")
		return nil
	}

	_, err := client.connection.pbClient.SendChat(context.Background(), &pb.SendChatRequest{
		Self:     client.user,
		Body:     s.Args.Message,
		Chatroom: client.subscription.chatroom,
	})

	return err
}
