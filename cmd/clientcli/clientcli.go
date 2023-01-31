package main

import (
	"log"

	"github.com/Richie78321/groupchat/client"
)

func main() {
	if err := client.Start("localhost:3000"); err != nil {
		log.Fatalf("%v", err)
	}
}
