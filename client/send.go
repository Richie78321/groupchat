package client

import (
	"context"

	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/buger/goterm"
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
	// We assume that having a current subscription implies having
	// a valid username and an active connection to the server.
	if client.subscription == nil {
		goterm.Println("Not in a chatroom")
		return nil
	}

	_, err := client.pbClient.SendChat(context.Background(), &pb.SendChatRequest{
		Self:     client.user,
		Body:     s.Args.Message,
		Chatroom: client.subscription.chatroom,
	})

	return err
}
