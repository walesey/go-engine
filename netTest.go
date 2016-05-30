package main

import (
	"github.com/walesey/go-engine/networking"
)

func main() {
	server := networking.NewServer()
	server.Listen(19999)
	for {

	}
}
