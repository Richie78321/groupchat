package client

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/google/shlex"
	"github.com/jessevdk/go-flags"
)

var client struct {
	printLock    sync.Mutex
	pbClient     pb.ChatServiceClient
	user         *pb.User
	subscription *subscription
}

func loggedIn() bool {
	return client.user != nil
}

func connected() bool {
	return client.pbClient != nil
}

func inChatroom() bool {
	return client.subscription != nil
}

var parser = flags.NewParser(&struct{}{}, flags.HelpFlag)

func splitArgs(text string) ([]string, error) {
	// Special-case arg splitting for message sending: everything
	// after the command is included in the message
	if strings.HasPrefix(text, "a ") {
		return []string{"a", text[2:]}, nil
	}

	return shlex.Split(text)
}

func printSeparator() {
	fmt.Printf("\n---\n")
}

func Start() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		splitArgs, err := splitArgs(scanner.Text())
		if err != nil {
			log.Fatalf("%v", err)
		}

		client.printLock.Lock()

		_, err = parser.ParseArgs(splitArgs)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}

		printSeparator()
		client.printLock.Unlock()
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("%v", err)
	}
}
