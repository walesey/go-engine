package networking

import (
	"fmt"
	"net"
)

type Server struct {
	sessions    []Session
	newSessions chan Session
	conn        *net.UDPConn
}

func NewServer() *Server {
	return &Server{
		sessions:    make([]Session, 0),
		newSessions: make(chan Session),
	}
}

func (s *Server) Listen(port int) {
	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%v", port))
	if err != nil {
		fmt.Println("Error Resolving udp address: ", err)
		return
	}

	s.conn, err = net.ListenUDP("udp", serverAddr)
	if err != nil {
		fmt.Println("Error listening on udp address: ", err)
		return
	}

	buf := make([]byte, 65500)
	go func() {
		for s.conn != nil {
			n, addr, err := s.conn.ReadFromUDP(buf)

			data := buf[0:n]
			fmt.Println("Received ", string(data), " from ", addr)

			if err != nil {
				fmt.Println("Error: ", err)
				break
			}
		}
	}()
}

func (s *Server) Close() {
	s.conn.Close()
	s.conn = nil
}
