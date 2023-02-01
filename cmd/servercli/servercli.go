package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Richie78321/groupchat/server"
)

func main() {
	address := flag.String("address", "localhost", "the server address")
	port := flag.Int("port", 3000, "the server port")
	flag.Parse()

	if err := server.Start(fmt.Sprintf("%s:%d", *address, *port)); err != nil {
		log.Fatalf("%v", err)
	}
}
