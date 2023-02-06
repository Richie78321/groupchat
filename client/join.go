package client

import (
	"context"
	"io"
	"log"
	"sort"
	"strings"

	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/buger/goterm"
)

type subscription struct {
	name   string
	stream pb.ChatService_SubscribeChatroomClient
	cancel context.CancelFunc
	ctx    context.Context
}

func (s *subscription) printUpdate(update *pb.ChatroomSubscriptionUpdate) {
	client.printLock.Lock()
	defer client.printLock.Unlock()
	goterm.Clear()
	goterm.MoveCursor(0, 0)

	goterm.Printf("Group: %s\n", s.name)

	// De-deuplicate the usernames for display
	usernameMap := make(map[string]struct{})
	for _, user := range update.Participants {
		usernameMap[user.Username] = struct{}{}
	}
	usernames := []string{}
	for username, _ := range usernameMap {
		usernames = append(usernames, username)
	}
	// Sort the usernames for consistency (key order is not guaranteed when
	// de-duplicating)
	sort.Strings(usernames)
	goterm.Printf("Participants: %s\n", strings.Join(usernames, ", "))

	goterm.Flush()
}

func (s *subscription) ingestUpdates() {
	for {
		update, err := s.stream.Recv()
		if err == io.EOF || s.ctx.Err() == context.Canceled {
			break
		}
		if err != nil {
			log.Fatalf("%v", err)
		}

		s.printUpdate(update)
	}
}

func endSubscription() {
	if client.subscription == nil {
		return
	}

	goterm.Printf("Ending existing subscription to `%s`\n", client.subscription.name)
	client.subscription.cancel()
	client.subscription = nil
}

type joinArgs struct {
	Args struct {
		Group string `description:"group name"`
	} `positional-args:"yes" required:"yes"`
}

func init() {
	parser.AddCommand("j", "join a chatroom", "", &joinArgs{})
}

func (j *joinArgs) Execute(args []string) error {
	if client.pbClient == nil {
		goterm.Println("Not connected to a server")
		return nil
	}

	if client.user == nil {
		goterm.Println("Need to login first")
		return nil
	}

	// Existing subscription ends when joining a new chatroom
	endSubscription()

	ctx, cancel := context.WithCancel(context.Background())
	stream, err := client.pbClient.SubscribeChatroom(ctx, &pb.SubscribeChatroomRequest{
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
		name:   j.Args.Group,
		stream: stream,
		cancel: cancel,
		ctx:    ctx,
	}

	// Spawn a thread to ingest the updates
	go client.subscription.ingestUpdates()

	return nil
}
