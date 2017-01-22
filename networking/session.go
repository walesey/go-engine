package networking

import (
	"bytes"
	"fmt"
	"net"
	"time"
)

var tokens int64

type Session struct {
	token        string
	addr         *net.UDPAddr
	idleTimer    time.Time
	packetBuffer *bytes.Buffer
}

func NewSession(addr *net.UDPAddr) *Session {
	return &Session{
		token:        generateToken(),
		addr:         addr,
		idleTimer:    time.Now(),
		packetBuffer: new(bytes.Buffer),
	}
}

func generateToken() string {
	tokens++
	return fmt.Sprintf("%v", tokens)
}
