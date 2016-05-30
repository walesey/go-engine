package networking

import (
	"fmt"
	"net"
)

var tokens int64

type Session struct {
	token string
	addr  net.Addr
}

func NewSession(addr net.Addr) *Session {
	return &Session{
		token: generateToken(),
		addr:  addr,
	}
}

func generateToken() string {
	tokens++
	return fmt.Sprintf("%v", tokens)
}

func (s *Session) WriteMessage(data []byte) {
	sessionAddr, err := net.ResolveUDPAddr("udp", s.addr.String())
	if err != nil {
		fmt.Println("Error resolving session udp address: ", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, sessionAddr)
	if err != nil {
		fmt.Println("Error connecting to udp session address: ", err)
		return
	}

	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("Error writing udp message to session address: ", err)
	}
}
