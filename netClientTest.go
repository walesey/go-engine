package main

import (
	"fmt"

	"github.com/walesey/go-engine/networking"
	"time"
)

func main() {
	client := networking.NewClient()
	client.Connect("127.0.0.1:29999")
	defer client.Close()

	for {
		next, ok := client.GetNextMessage()
		if ok {
			fmt.Println(next)
		}
		client.WriteMessage("test-message", "arg1", "arg2", 123)
		time.Sleep(1000 * time.Millisecond)
	}
}
