package networking

import (
	"fmt"
	"net"
)

const serverPacketBufferSize = 1000
const serverSessionBufferSize = 20

type Server struct {
	sessions    []*Session
	newSessions chan *Session
	packets     chan Packet
}

func NewServer() *Server {
	return &Server{
		sessions:    make([]*Session, 0),
		newSessions: make(chan *Session, serverSessionBufferSize),
		packets:     make(chan Packet, clientPacketBufferSize),
	}
}

func (s *Server) Listen(port int) {
	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%v", port))
	if err != nil {
		fmt.Println("Error Resolving udp address: ", err)
		return
	}

	go func() {
		for {
			fmt.Println("test")
			newConn, err := net.ListenUDP("udp", serverAddr)
			if err != nil {
				fmt.Println("Error listening on udp address: ", err)
				break
			}
			s.newSessions <- NewSession(newConn, s.packets)
		}
	}()
}

func (s *Server) GetNextMessage() (Packet, bool) {
	select {
	case packet := <-s.packets:
		return packet, true
	default:
	}
	return Packet{}, false
}

func (s *Server) Update(dt float64) {
	select {
	case packet := <-s.newSessions:
		s.sessions = append(s.sessions, packet)
		break
	default:
	}
}

func (s *Server) CloseAll() {
	for _, session := range s.sessions {
		session.Close()
	}
}
