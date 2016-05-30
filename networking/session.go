package networking

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
)

var ids int64

type Command struct {
	command string
	args    []interface{}
}

type Packet struct {
	Command
	sessionId int64
}

type Session struct {
	id      int64
	conn    *net.UDPConn
	packets chan Packet
}

func NewSession(conn *net.UDPConn, packets chan Packet) *Session {
	ids++
	return &Session{
		id:      ids,
		conn:    conn,
		packets: packets,
	}
}

func (s *Session) Listen() {
	data := make([]byte, 65500)
	go func() {
		for s.conn != nil {
			n, _, err := s.conn.ReadFromUDP(data)
			if err != nil {
				fmt.Println("Error reading udp packet: ", err)
				continue
			}

			dataBuf := bytes.NewBuffer(data[0:n])
			decoder := gob.NewDecoder(dataBuf)
			var command Command
			err = decoder.Decode(&command)
			if err != nil {
				fmt.Println("Error decoding udp packet: ", err)
				continue
			}

			s.packets <- Packet{Command: command, sessionId: s.id}
		}
	}()
}

func (s *Session) Close() {
	s.conn.Close()
	s.conn = nil
}

func (s *Session) WriteMessage(data []byte) {
	_, err := s.conn.Write(data)
	if err != nil {
		fmt.Println("Error writing udp message to session address: ", err)
	}
}
