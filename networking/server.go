package networking

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net"
	"time"
)

const sessionTimeout = 10 * time.Minute

type Server struct {
	conn             *net.UDPConn
	sessions         map[string]*Session
	onClientJoined   func(clientId string)
	onPacketReceived func(packet Packet)
	bytesSent        int64
	bytesReceived    int64
}

func NewServer() *Server {
	var server *Server
	server = &Server{
		sessions: make(map[string]*Session),
		onClientJoined: func(clientId string) {
			server.WriteMessage("", clientId, []byte{})
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

			s.bytesReceived += int64(n)
			dataBuf := bytes.NewBuffer(data[0:n])
			gzipReader, err := gzip.NewReader(dataBuf)
			if err != nil {
				fmt.Println("Error creating gzip Reader for udp packet: ", err)
				continue
			}

			unzipped, err := ioutil.ReadAll(gzipReader)
			if err != nil {
				fmt.Println("Error unzipping udp packet: ", err)
				continue
			}

			var packet Packet
			for i := 0; i < len(unzipped); {
				packet, err, i = Decode(unzipped, i)
				if err != nil {
					fmt.Println("Error Decoding udp packet: ", err)
					continue
				}

				if len(packet.Token) == 0 {
					newSession := NewSession(addr)
					s.sessions[newSession.token] = newSession
					s.onClientJoined(newSession.token)
					s.cleanupSessions()
				}

				if session, ok := s.sessions[packet.Token]; ok {
					session.idleTimer = time.Now()
				}

				if s.onPacketReceived != nil {
					s.onPacketReceived(packet)
				}
			}
		}
	}()
}

func (s *Server) WriteMessage(command, token string, data []byte) {
	if session, ok := s.sessions[token]; ok {
		packet := Packet{
			Token:   token,
			Command: command,
			Data:    data,
		}

		if _, err := session.packetBuffer.Write(Encode(packet)); err != nil {
			fmt.Println("Error Writing udp message to session buffer: ", err)
		}
	}
}

func (s *Server) cleanupSessions() {
	// check for session timeouts
	for token, session := range s.sessions {
		if time.Since(session.idleTimer) > sessionTimeout {
			fmt.Println("session timed out: ", token)
			delete(s.sessions, token)
		}
		session.idleTimer = time.Now()
	}
}

func (s *Server) BroadcastMessage(command string, data []byte) {
	for token, _ := range s.sessions {
		s.WriteMessage(command, token, data)
	}
}

func (s *Server) ClientJoinedEvent(callback func(clientId string)) {
	s.onClientJoined = callback
}

func (s *Server) PacketReceived(callback func(packet Packet)) {
	s.onPacketReceived = callback
}

// FlushAllWriteBuffers - send all buffered messages immediately for all sessions
func (s *Server) FlushAllWriteBuffers() {
	for token := range s.sessions {
		s.FlushWriteBuffer(token)
	}
}

// FlushWriteBuffer - send all buffered messages immediately
func (s *Server) FlushWriteBuffer(token string) {
	if session, ok := s.sessions[token]; ok {
		data := session.packetBuffer.Bytes()
		if len(data) > 0 {
			var gzipBuf bytes.Buffer
			gzipWriter := gzip.NewWriter(&gzipBuf)
			_, err := gzipWriter.Write(data)
			if err != nil {
				fmt.Println("Error Gzip compressing udp message: ", err)
				return
			}

			if err := gzipWriter.Flush(); err != nil {
				fmt.Println("Error Flushing Gzip writer for udp message: ", err)
				return
			}

			if err := gzipWriter.Close(); err != nil {
				fmt.Println("Error Closing Gzip writer for udp message: ", err)
				return
			}

			gzipData := gzipBuf.Bytes()
			s.bytesSent += int64(len(gzipData))
			s.conn.WriteToUDP(gzipData, session.addr)
			session.packetBuffer.Reset()
		}
	}
}

func (s *Server) Close() {
	s.FlushAllWriteBuffers()
	s.conn.Close()
	s.conn = nil
}
