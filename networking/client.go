package networking

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net"
)

const clientPacketBufferSize = 100

type Client struct {
	token                string
	conn                 *net.UDPConn
	packets              chan Packet
	bytesSent            int64
	bytesReceived        int64
	bytesSentByEvent     map[string]int64
	bytesReceivedByEvent map[string]int64
}

func NewClient() *Client {
	return &Client{
		packets:              make(chan Packet, clientPacketBufferSize),
		bytesSentByEvent:     make(map[string]int64),
		bytesReceivedByEvent: make(map[string]int64),
	}
}

func (c *Client) Connect(addr string) {
	serverAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Println("Error resolving server udp address: ", err)
		return
	}

	c.conn, err = net.DialUDP("udp", nil, serverAddr)
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
			gzipReader, err := gzip.NewReader(dataBuf)
			if err != nil {
				fmt.Println("Error creating gzip Reader for udp packet: ", err)
				continue
			}

			unzipped, err := ioutil.ReadAll(gzipReader)
			if err != nil {
				fmt.Println("Error unzipping udp packet: ", err)
				continue
			}

			var packet Packet
			for i := 0; i < len(unzipped); {
				j := i
				packet, err, i = Decode(unzipped, i)
				if err != nil {
					fmt.Println("Error decoding udp packet: ", err)
					continue
				}
				c.updateBytesReceived(packet.Command, int64(i-j))

				c.token = packet.Token
				c.packets <- packet
			}
		}
	}()
}

func (c *Client) WriteMessage(command string, data []byte) {
	packet := Packet{
		Token:   c.token,
		Command: command,
		Data:    data,
	}

	packetData := Encode(packet)
	var gzipBuf bytes.Buffer
	gzipWriter := gzip.NewWriter(&gzipBuf)
	_, err := gzipWriter.Write(packetData)
	if err != nil {
		fmt.Println("Error Gzip compressing udp message: ", err)
		return
	}

	if err := gzipWriter.Flush(); err != nil {
		fmt.Println("Error Flushing Gzip writer for udp message: ", err)
		return
	}

	if err := gzipWriter.Close(); err != nil {
		fmt.Println("Error Closing Gzip writer for udp message: ", err)
		return
	}

	gzipData := gzipBuf.Bytes()
	c.updateBytesSent(command, int64(len(gzipData)))
	_, err = c.conn.Write(gzipData)
	if err != nil {
		fmt.Println("Error writing udp message: ", err)
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

func (c *Client) updateBytesSent(event string, sent int64) {
	c.bytesSent += sent
	total, ok := c.bytesSentByEvent[event]
	if !ok {
		c.bytesSentByEvent[event], total = 0, 0
	}
	c.bytesSentByEvent[event] = sent + total
}

func (c *Client) updateBytesReceived(event string, sent int64) {
	c.bytesReceived += sent
	total, ok := c.bytesReceivedByEvent[event]
	if !ok {
		c.bytesReceivedByEvent[event], total = 0, 0
	}
	c.bytesReceivedByEvent[event] = sent + total
}
