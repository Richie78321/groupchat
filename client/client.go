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
	"google.golang.org/grpc"
)

var client struct {
	printLock    sync.Mutex
	user         *pb.User
	subscription *subscription
	connection   struct {
		grpc     *grpc.ClientConn
		pbClient pb.ChatServiceClient
	}

	shouldExit bool
}

func loggedIn() bool {
	return client.user != nil
}

func connected() bool {
	return client.connection.pbClient != nil
}

// inChatroom subsumes loggedIn() and connected(), as those are prerequisites
// for being subscribed to a chatroom.
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
	client.shouldExit = false

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		splitArgs, err := splitArgs(scanner.Text())
		if err != nil {
			log.Fatalf("%v", err)
		}

		// Before running a command, obtain the printLock to avoid
		// concurrent writes to the console.
		client.printLock.Lock()

		// ParseArgs parses the command and executes the respective Execute() handler.
		// Any errors returned from handlers are returned here.
		_, err = parser.ParseArgs(splitArgs)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}

		printSeparator()
		client.printLock.Unlock()

		if client.shouldExit {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("%v", err)
	}

	// Clean up any open connections
	if client.connection.grpc != nil {
		client.connection.grpc.Close()
	}
}
