package client

import (
	"context"
	"fmt"

	pb "github.com/Richie78321/groupchat/chatservice"
)

type joinArgs struct {
	Args struct {
		Group string `description:"group name"`
	} `positional-args:"yes" required:"yes"`
}

func init() {
	parser.AddCommand("j", "join a chatroom", "", &joinArgs{})
}

func (j *joinArgs) Execute(args []string) error {
	if !connected() {
		fmt.Println("Not connected to a server")
		return nil
	}

	if !loggedIn() {
		fmt.Println("Need to login first")
		return nil
	}

	// Existing subscription ends when joining a new chatroom
	endSubscription()

	// Use a cancellable context so the client can terminate subscriptions.
	ctx, cancel := context.WithCancel(context.Background())
	stream, err := client.connection.pbClient.SubscribeChatroom(ctx, &pb.SubscribeChatroomRequest{
		Self: client.user,
		Chatroom: &pb.Chatroom{
			Name: j.Args.Group,
		},
	})
	if err != nil {
		cancel()
		return err
	}

	client.subscription = &subscription{
		chatroom: &pb.Chatroom{
			Name: j.Args.Group,
		},
		stream: stream,
		cancel: cancel,
		ctx:    ctx,
	}

	// Spawn a thread to ingest subscription updates
	go client.subscription.ingestUpdates()

	return nil
}
