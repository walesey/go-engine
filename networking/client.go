package networking

import (
	"fmt"
	"net"
)

const clientPacketBufferSize = 100

type Client struct {
	session *Session
	packets chan Packet
}

func NewClient() *Client {
	return &Client{
		packets: make(chan Packet, clientPacketBufferSize),
	}
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

	conn, err := net.DialUDP("udp", localAddr, serverAddr)
	if err != nil {
		fmt.Println("Error connecting to udp server address: ", err)
		return
	}

	c.session = NewSession(conn, c.packets)
	c.session.Listen()
}

func (c *Client) WriteMessage(data []byte) {
	c.session.WriteMessage(data)
}

func (c *Client) GetNextMessage() (Packet, bool) {
	select {
	case packet := <-c.packets:
		return packet, true
	default:
	}
	return Packet{}, false
}

func (c *Client) Close() {
	c.session.Close()
}
