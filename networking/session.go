package networking

import (
	"fmt"
	"net"
)

var tokens int64

type Session struct {
	token string
	addr  *net.UDPAddr
}

func NewSession(addr *net.UDPAddr) *Session {
	return &Session{
		token: generateToken(),
		addr:  addr,
	}
}

func generateToken() string {
	tokens++
	return fmt.Sprintf("%v", tokens)
}
