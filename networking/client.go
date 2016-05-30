package networking

import (
	"fmt"
	"net"
)

type Client struct {
	session Session
	conn    *net.UDPConn
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Connect(addr string) {
	serverAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Println("Error resolving server udp address: ", err)
		return
	}

	localAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("Error resolving local udp address: ", err)
		return
	}

	c.conn, err = net.DialUDP("udp", localAddr, serverAddr)
	if err != nil {
		fmt.Println("Error connecting to udp server address: ", err)
		return
	}

	buf := make([]byte, 65500)
	go func() {
		for c.conn != nil {
			n, addr, err := c.conn.ReadFromUDP(buf)
			data := buf[0:n]
			fmt.Println("Received ", string(data), " from ", addr)

			if err != nil {
				fmt.Println("Error: ", err)
				break
			}
		}
	}()
}

func (c *Client) WriteMessage(data []byte) {
	_, err := c.conn.Write(data)
	if err != nil {
		fmt.Println("Error writting message to server: ", err)
		return
	}
}

func (c *Client) Close() {
	c.conn.Close()
	c.conn = nil
}
