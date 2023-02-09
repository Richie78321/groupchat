package client

import (
	"context"
	"fmt"

	pb "github.com/Richie78321/groupchat/chatservice"
)

type historyArgs struct{}

func init() {
	parser.AddCommand("p", "get full message history", "", &historyArgs{})
}

func (h *historyArgs) Execute(args []string) error {
	if !inChatroom() {
		fmt.Println("Not in chatroom")
		return nil
	}

	response, err := client.pbClient.MessageHistory(context.Background(), &pb.MessageHistoryRequest{
		Chatroom: client.subscription.chatroom,
	})
	if err != nil {
		return err
	}

	client.subscription.showLatestMessages(response.Messages)
	return nil
}
