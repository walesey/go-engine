package networking

import (
	"fmt"
	"net"
)

var tokenId = 0

type Session struct {
	token string
	addr  net.Addr
}

func NewSession(addr net.Addr) Session {
	return Session{addr: addr, token: generateToken()}
}

func generateToken() string {
	tokenId++
	return fmt.Sprintf("%v", tokenId)
}

func (sess Session) WriteMessage(data []byte) {
	clientAddr, err := net.ResolveUDPAddr("udp", sess.addr.String())
	if err != nil {
		fmt.Println("Error resolving server udp address: ", err)
		return
	}

	localAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("Error resolving local udp address: ", err)
		return
	}

	clientConn, err := net.DialUDP("udp", localAddr, clientAddr)
	if err != nil {
		fmt.Println("Error connecting to udp client address: ", err)
		return
	}

	_, err = clientConn.Write(data)
	if err != nil {
		fmt.Println("Error writting message to client: ", err)
		return
	}
}
