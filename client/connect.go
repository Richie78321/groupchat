package client

import (
	"fmt"
	"time"

	pb "github.com/Richie78321/groupchat/chatservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type connectArgs struct {
	Args struct {
		Address string `description:"server address"`
	} `positional-args:"yes" required:"yes"`
}

func init() {
	parser.AddCommand("c", "connect to a server", "", &connectArgs{})
}

func (c *connectArgs) Execute(args []string) error {
	if connected() {
		fmt.Println("Already connected")
		return nil
	}

	fmt.Println("Connecting...")
	conn, err := grpc.Dial(c.Args.Address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock(), grpc.WithTimeout(10*time.Second))
	if err != nil {
		return err
	}

	client.connection.grpc = conn
	client.connection.pbClient = pb.NewChatServiceClient(conn)
	fmt.Printf("Connected to `%s`\n", c.Args.Address)
	return nil
}
