package networking

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
)

const clientPacketBufferSize = 100

type Client struct {
	token   string
	conn    *net.UDPConn
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

	c.conn, err = net.DialUDP("udp", localAddr, serverAddr)
	if err != nil {
		fmt.Println("Error connecting to udp server address: ", err)
		return
	}

	data := make([]byte, 65500)
	go func() {
		for c.conn != nil {
			n, _, err := c.conn.ReadFromUDP(data)
			if err != nil {
				fmt.Println("Error reading udp packet: ", err)
				continue
			}

			dataBuf := bytes.NewBuffer(data[0:n])
			decoder := gob.NewDecoder(dataBuf)
			var packet Packet
			err = decoder.Decode(&packet)
			if err != nil {
				fmt.Println("Error decoding udp packet: ", err)
				continue
			}
			c.token = packet.Token

			c.packets <- packet
		}
	}()
}

func (c *Client) WriteMessage(command string, args ...interface{}) {
	packet := Packet{
		Token:   c.token,
		Command: command,
		Args:    args,
	}

	dataBuf := bytes.NewBuffer(make([]byte, 61500))
	endcoder := gob.NewEncoder(dataBuf)
	err := endcoder.Encode(packet)
	if err != nil {
		fmt.Println("Error encoding udp message: ", err)
		return
	}

	_, err = c.conn.Write(dataBuf.Bytes())
	if err != nil {
		fmt.Println("Error writing udp message to session address: ", err)
	}
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
	c.conn.Close()
}
