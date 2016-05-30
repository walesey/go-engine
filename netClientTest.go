package main

import (
	"fmt"

	"github.com/walesey/go-engine/networking"
)

func main() {
	client := networking.NewClient()
	client.Connect("127.0.0.1:29999")
	defer client.Close()

	client.WriteMessage("test client message")
	for {
		next, ok := client.GetNextMessage()
		if ok {
			fmt.Println(next)
		}
	}
}
