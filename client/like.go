package client

import (
	"context"
	"fmt"

	pb "github.com/Richie78321/groupchat/chatservice"
)

func likeOrUnlike(messageId int, like bool) error {
	if !inChatroom() {
		fmt.Println("Not in a chatroom")
		return nil
	}

	if client.subscription.latestMessages == nil {
		fmt.Println("There are currently no messages to select from")
		return nil
	}

	shiftedId := messageId - 1
	if shiftedId < 0 || shiftedId >= len(client.subscription.latestMessages) {
		fmt.Printf("The specified ID `%d` is out of range\n", messageId)
		return nil
	}

	_, err := client.connection.pbClient.LikeChat(context.Background(), &pb.LikeChatRequest{
		Self:        client.user,
		MessageUuid: client.subscription.latestMessages[shiftedId].Uuid,
		Chatroom:    client.subscription.chatroom,
		Like:        like,
	})

	return err
}

func init() {
	parser.AddCommand("l", "like a message", "", &likeArgs{})
	parser.AddCommand("r", "unlike a message", "", &unlikeArgs{})
}

type likeArgs struct {
	Args struct {
		MessageId int `description:"the displayed ID of the message"`
	} `positional-args:"yes" required:"yes"`
}

func (l *likeArgs) Execute(args []string) error {
	return likeOrUnlike(l.Args.MessageId, true)
}

type unlikeArgs struct {
	Args struct {
		MessageId int `description:"the displayed ID of the message"`
	} `positional-args:"yes" required:"yes"`
}

func (u *unlikeArgs) Execute(args []string) error {
	return likeOrUnlike(u.Args.MessageId, false)
}
