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
	chatroom *pb.Chatroom
	stream   pb.ChatService_SubscribeChatroomClient
	cancel   context.CancelFunc
	ctx      context.Context
}

func (s *subscription) printUpdate(update *pb.ChatroomSubscriptionUpdate) {
	client.printLock.Lock()
	defer client.printLock.Unlock()
	goterm.Clear()
	goterm.MoveCursor(0, 0)

	goterm.Printf("Group: %s\n", s.chatroom.Name)

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

	for i, m := range update.LatestMessages {
		// TODO(richie): Add likers when they are implemented
		goterm.Printf("%d. %s: %s\n", i+1, m.Author.Username, m.Body)
	}

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

	goterm.Printf("Ending existing subscription to `%s`\n", client.subscription.chatroom.Name)
	client.subscription.cancel()
	client.subscription = nil
}
