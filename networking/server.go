package networking

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"

	"github.com/walesey/go-engine/util"
	vmath "github.com/walesey/go-engine/vectormath"
)

const serverPacketBufferSize = 1000
const serverSessionBufferSize = 20

func init() {
	gob.Register([]interface{}{})
	gob.Register(vmath.Vector4{})
	gob.Register(vmath.Vector3{})
	gob.Register(vmath.Vector2{})
	gob.Register(vmath.Quaternion{})
}

type Packet struct {
	Token   string
	Command string
	Args    []interface{}
}

type Server struct {
	conn           *net.UDPConn
	sessions       map[string]*Session
	newSessions    chan *Session
	packets        chan Packet
	onClientJoined func(clientId string)
}

func NewServer() *Server {
	var server *Server
	server = &Server{
		sessions:    make(map[string]*Session),
		newSessions: make(chan *Session, serverSessionBufferSize),
		packets:     make(chan Packet, clientPacketBufferSize),
		onClientJoined: func(clientId string) {
			server.WriteMessage("", clientId)
		},
	}
	return server
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

	data := make([]byte, 65500)
	go func() {
		for s.conn != nil {
			n, addr, err := s.conn.ReadFromUDP(data)
			if err != nil {
				fmt.Println("Error reading udp packet: ", err)
				continue
			}

			dataBuf := bytes.NewBuffer(data[0:n])
			decoder := gob.NewDecoder(dataBuf)
			var packet Packet
			err = decoder.Decode(&packet)
			if err != nil {
				fmt.Println("Error decoding udp packet: ", err)
				continue
			}

			if len(packet.Token) == 0 {
				s.newSessions <- NewSession(addr)
			}

			s.packets <- packet
		}
	}()
}

func (s *Server) WriteMessage(command string, token string, args ...interface{}) {
	session, ok := s.sessions[token]
	if ok {
		packet := Packet{
			Token:   token,
			Command: command,
			Args:    args,
		}

		data, err := util.Serialize(packet)
		if err != nil {
			fmt.Println("Error Serializing udp message: ", err)
			return
		}

		s.conn.WriteToUDP(data, session.addr)
		if err != nil {
			fmt.Println("Error Writing udp message: ", err)
		}
	}
}

func (s *Server) BroadcastMessage(command string, args ...interface{}) {
	for token, _ := range s.sessions {
		s.WriteMessage(command, token, args...)
	}
}

func (s *Server) ClientJoinedEvent(callback func(clientId string)) {
	s.onClientJoined = callback
}

func (s *Server) Update(dt float64) {
	select {
	case newSession := <-s.newSessions:
		s.sessions[newSession.token] = newSession
		s.onClientJoined(newSession.token)
	default:
	}
}

func (s *Server) GetNextMessage() (Packet, bool) {
	select {
	case packet := <-s.packets:
		return packet, true
	default:
	}
	return Packet{}, false
}

func (s *Server) Close() {
	s.conn.Close()
	s.conn = nil
}
