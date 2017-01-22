package networking

import (
	"time"

	"github.com/walesey/go-engine/emitter"
	"github.com/walesey/go-engine/util"
)

type Network struct {
	*emitter.Emitter
	stopInterval func()
	client       *Client
	server       *Server
}

func NewNetwork() *Network {
	var network *Network
	network = &Network{
		Emitter: emitter.New(16),
	}
	return network
}

func (n *Network) StartServer(port int) {
	n.server = NewServer()
	n.server.ClientJoinedEvent(func(clientId string) {
		n.Emit("newClient", clientId)
	})
	n.server.PacketReceived(func(packet Packet) {
		n.Emit(packet.Command, packet)
	})
	n.server.Listen(port)
	n.startWriteInterval(50 * time.Millisecond)
}

func (n *Network) ConnectClient(addr string) error {
	n.client = NewClient()
	n.client.PacketReceived(func(packet Packet) {
		n.Emit(packet.Command, packet)
	})
	if err := n.client.Connect(addr); err != nil {
		return err
	}
	n.client.WriteMessage("", []byte{})
	return nil
}

// Update is used when using network in a fully syncronous manner.
// Update should not be called when using asyncronous (channel) based event handling
func (n *Network) Update(dt float64) {
	n.FlushAll()
}

func (n *Network) ClientJoinedEvent(fn func(clientId string)) {
	n.On("newClient", func(event emitter.Event) {
		if clientId, ok := event.(string); ok {
			fn(clientId)
		}
	})
}

// RegisterEvent - register an event that will be triggered on clients and server.
func (n *Network) RegisterEvent(name string, fn func(clientId string, data []byte)) {
	n.On(name, func(event emitter.Event) {
		if packet, ok := event.(Packet); ok {
			fn(packet.Token, packet.Data)
		}
	})
}

// TriggerEvent - Trigger an event to run on a particular client.
// If called on the client, this will trigger the event on the server.
func (n *Network) TriggerEvent(name, clientId string, data []byte) {
	if n.IsClient() {
		n.client.WriteMessage(name, data)
	}
	if n.IsServer() {
		n.server.WriteMessage(name, clientId, data)
	}
}

// BroadcastEvent - trigger an event on all clients.
// If called on the client, this will trigger the event on the server.
func (n *Network) BroadcastEvent(name string, data []byte) {
	if n.IsClient() {
		n.client.WriteMessage(name, data)
	}
	if n.IsServer() {
		n.server.BroadcastMessage(name, data)
	}
}

// CallOnServerAndClient - trigger an event on the server and on all client.
// If called on the client, this will trigger the event on the client and on the server.
func (n *Network) TriggerOnServerAndClients(name string, data []byte) {
	n.Emit(name, Packet{Data: data})
	if n.IsClient() {
		n.client.WriteMessage(name, data)
	}
	if n.IsServer() {
		n.server.BroadcastMessage(name, data)
	}
}

func (n *Network) ClientToken() string {
	return n.client.token
}

func (n *Network) FlushAllWriteBuffers() {
	if n.IsServer() {
		n.server.FlushAllWriteBuffers()
	}
}

func (n *Network) FlushWriteBuffer(clientId string) {
	if n.IsServer() {
		n.server.FlushWriteBuffer(clientId)
	}
}

func (n *Network) BytesSent() int64 {
	if n.IsClient() {
		return n.client.bytesSent
	}
	if n.IsServer() {
		return n.server.bytesSent
	}
	return 0
}

func (n *Network) BytesReceived() int64 {
	if n.IsClient() {
		return n.client.bytesReceived
	}
	if n.IsServer() {
		return n.server.bytesReceived
	}
	return 0
}

func (n *Network) BytesSentByCommand() map[string]int64 {
	if n.IsClient() {
		return n.client.bytesSentByEvent
	}
	return map[string]int64{}
}

func (n *Network) BytesReceivedByCommand() map[string]int64 {
	if n.IsClient() {
		return n.client.bytesReceivedByEvent
	}
	return map[string]int64{}
}

func (n *Network) IsClient() bool {
	return n.client != nil
}

func (n *Network) IsServer() bool {
	return n.server != nil
}

func (n *Network) Kill() {
	n.KillClient()
	n.KillServer()
	n.Close()
}

func (n *Network) KillClient() {
	if n.client != nil {
		n.client.Close()
	}
	n.stopIntervals()
	n.client = nil
}

func (n *Network) KillServer() {
	if n.server != nil {
		n.server.Close()
		n.Off("newClient")
	}
	n.stopIntervals()
	n.server = nil
}

func (n *Network) startWriteInterval(bufferedWriteDuration time.Duration) {
	n.stopIntervals()
	if n.server != nil {
		n.stopInterval = util.SetInterval(n.server.FlushAllWriteBuffers, bufferedWriteDuration)
	}
}

func (n *Network) stopIntervals() {
	if n.stopInterval != nil {
		n.stopInterval()
	}
}
