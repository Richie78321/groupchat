package main

import (
	"log"

	"github.com/Richie78321/groupchat/server"
)

func main() {
	if err := server.Start("localhost:3000"); err != nil {
		log.Fatalf("%v", err)
	}
}
