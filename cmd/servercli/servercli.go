package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/Richie78321/groupchat/server"
	"github.com/Richie78321/groupchat/server/replication"
	"github.com/jessevdk/go-flags"
)

var opts struct {
	Address string   `long:"address" description:"The server address" default:"localhost"`
	Port    int      `long:"port" description:"The server port" default:"3000"`
	Peers   []string `long:"peer" short:"p" description:"Peer server in the format <id>:<address>. This can be called multiple times."`
	Args    struct {
		Id string `description:"the server's ID"`
	} `positional-args:"yes" required:"yes"`
}

func peersFromArgs() ([]*replication.Peer, error) {
	peers := make([]*replication.Peer, 0, len(opts.Peers))
	for _, peer := range opts.Peers {
		split := strings.SplitN(peer, ":", 2)
		if len(split) != 2 {
			return nil, fmt.Errorf("invalid peer string `%s`", peer)
		}

		peers = append(peers, replication.NewPeer(split[0], split[1]))
	}

	return peers, nil
}

func main() {
	parser := flags.NewParser(&opts, flags.HelpFlag)
	_, err := parser.Parse()
	if err != nil {
		log.Fatalf("%v", err)
	}

	peers, err := peersFromArgs()
	if err != nil {
		log.Fatalf("%v", err)
	}

	if err := server.Start(opts.Args.Id, fmt.Sprintf("%s:%d", opts.Address, opts.Port), peers); err != nil {
		log.Fatalf("%v", err)
	}
}
