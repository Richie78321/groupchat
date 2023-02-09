package client

import (
	"bufio"
	"log"
	"os"
	"strings"
	"sync"

	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/buger/goterm"
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

func Start() {
	// Do initial screen clear
	client.printLock.Lock()
	goterm.Clear()
	goterm.MoveCursor(0, 0)
	goterm.Flush()
	client.printLock.Unlock()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		splitArgs, err := splitArgs(scanner.Text())
		if err != nil {
			log.Fatalf("%v", err)
		}

		client.printLock.Lock()
		goterm.Clear()
		goterm.MoveCursor(0, 0)

		_, err = parser.ParseArgs(splitArgs)
		if err != nil {
			goterm.Printf("error: %v\n", err)
		}

		goterm.Flush()
		client.printLock.Unlock()
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("%v", err)
	}
}
