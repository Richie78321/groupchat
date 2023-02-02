package client

import (
	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/buger/goterm"
)

type loginArgs struct {
	Args struct {
		Username string `description:"username"`
	} `positional-args:"yes" required:"yes"`
}

func init() {
	parser.AddCommand("u", "login", "", &loginArgs{})
}

func (j *loginArgs) Execute(args []string) error {
	// Existing subscription ends when logging in as a new user
	endSubscription()

	client.user = &pb.User{
		Username: j.Args.Username,
	}

	goterm.Printf("Logged in as user `%s`\n", client.user.Username)
	return nil
}
