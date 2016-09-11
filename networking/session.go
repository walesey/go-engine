package networking

import (
	"bytes"
	"fmt"
	"net"
)

var tokens int64

type Session struct {
	token        string
	addr         *net.UDPAddr
	idleTime     float64
	packetBuffer *bytes.Buffer
}

func NewSession(addr *net.UDPAddr) *Session {
	return &Session{
		token:        generateToken(),
		addr:         addr,
		idleTime:     0,
		packetBuffer: new(bytes.Buffer),
	}
}

func generateToken() string {
	tokens++
	return fmt.Sprintf("%v", tokens)
}
