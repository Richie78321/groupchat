package client

import (
	"context"
	"fmt"
	"strings"

	pb "github.com/Richie78321/groupchat/chatservice"
)

type viewPeersArgs struct{}

func init() {
	parser.AddCommand("v", "view peer servers that this server is connected to", "", &viewPeersArgs{})
}

func (v *viewPeersArgs) Execute(args []string) error {
	if !connected() {
		fmt.Println("Not connected")
		return nil
	}

	response, err := client.connection.pbClient.ViewPeers(context.Background(), &pb.ViewPeersRequest{})
	if err != nil {
		return err
	}

	if len(response.Peers) <= 0 {
		fmt.Println("Server is connected to zero peers")
		return nil
	}

	peerIds := make([]string, 0, len(response.Peers))
	for _, peer := range response.Peers {
		peerIds = append(peerIds, peer.Id)
	}
	fmt.Printf("Server is connected to peers: %s", strings.Join(peerIds, ", "))
	return nil
}
