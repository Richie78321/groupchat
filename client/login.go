package client

import (
	"fmt"

	pb "github.com/Richie78321/groupchat/chatservice"
)

type loginArgs struct {
	Args struct {
		Username string `description:"username"`
	} `positional-args:"yes" required:"yes"`
}

func init() {
	parser.AddCommand("u", "login", "", &loginArgs{})
}

func (l *loginArgs) Execute(args []string) error {
	// Existing subscription ends when logging in as a new user
	endSubscription()

	client.user = &pb.User{
		Username: l.Args.Username,
	}

	fmt.Printf("Logged in as user `%s`\n", client.user.Username)
	return nil
}
