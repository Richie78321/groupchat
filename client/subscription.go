package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"sort"
	"strings"
	"sync"

	pb "github.com/Richie78321/groupchat/chatservice"
)

type subscription struct {
	chatroom *pb.Chatroom
	stream   pb.ChatService_SubscribeChatroomClient
	cancel   context.CancelFunc
	ctx      context.Context

	// The subscription lock should be strictly acquired after acquiring the print lock.
	// Otherwise deadlocks are possible.
	lock           sync.Mutex
	latestMessages []*pb.Message
}

func (s *subscription) showLatestMessages(messages []*pb.Message) {
	s.lock.Lock()
	defer s.lock.Unlock()

	// Update the latest messages in the subscription so other functions
	// can reference them (for example to like a message by display index)
	s.latestMessages = messages

	if len(s.latestMessages) <= 0 {
		fmt.Printf("<No messages to display>\n")
		return
	}

	for i, m := range s.latestMessages {
		likersStr := ""
		if len(m.Likers) > 0 {
			likersStr = fmt.Sprintf(" [Likers: %d]", len(m.Likers))
		}

		fmt.Printf("%d. %s: %s%s\n", i+1, m.Author.Username, m.Body, likersStr)
	}
}

func (s *subscription) printUpdate(update *pb.ChatroomSubscriptionUpdate) {
	client.printLock.Lock()
	defer client.printLock.Unlock()

	fmt.Printf("Group: %s\n", s.chatroom.Name)

	// De-deuplicate the usernames for display
	usernameMap := make(map[string]struct{})
	for _, user := range update.Participants {
		usernameMap[user.Username] = struct{}{}
	}
	usernames := []string{}
	for username := range usernameMap {
		usernames = append(usernames, username)
	}
	// Sort the usernames for consistency (key order is not guaranteed when
	// de-duplicating)
	sort.Strings(usernames)
	fmt.Printf("Participants: %s\n", strings.Join(usernames, ", "))

	s.showLatestMessages(update.LatestMessages)

	printSeparator()
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

	fmt.Printf("Ending existing subscription to `%s`\n", client.subscription.chatroom.Name)
	client.subscription.cancel()
	client.subscription = nil
}
