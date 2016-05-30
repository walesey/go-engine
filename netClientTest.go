package main

import (
	"github.com/walesey/go-engine/networking"
)

func main() {
	client := networking.NewClient()
	client.Connect("127.0.0.1:19999")
	for {
	}
}
