package client

import (
	"bufio"
	"log"
	"os"
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

var parser = flags.NewParser(&struct{}{}, flags.HelpFlag)

func Start() {
	// Do initial screen clear
	client.printLock.Lock()
	goterm.Clear()
	goterm.MoveCursor(0, 0)
	goterm.Flush()
	client.printLock.Unlock()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		splitArgs, err := shlex.Split(scanner.Text())
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
