package main

import (
	"fmt"
	"time"

	"github.com/walesey/go-engine/networking"
)

func main() {
	server := networking.NewServer()
	server.Listen(29999)
	defer server.Close()
	for {
		server.Update(0)
		next, ok := server.GetNextMessage()
		if ok {
			fmt.Println(next)
		}

		time.Sleep(1000 * time.Millisecond)
	}
}
