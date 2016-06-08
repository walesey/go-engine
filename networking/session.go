package networking

import (
	"fmt"
	"net"
)

var tokens int64

type Session struct {
	token    string
	addr     *net.UDPAddr
	idleTime float64
}

func NewSession(addr *net.UDPAddr) *Session {
	return &Session{
		token:    generateToken(),
		addr:     addr,
		idleTime: 0,
	}
}

func generateToken() string {
	tokens++
	return fmt.Sprintf("%v", tokens)
}
