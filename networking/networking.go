package networking

import "log"

type Network struct {
	client *Client
	server *Server
	events map[string]func(clientId string, data []byte)
}

func NewNetwork() *Network {
	var network *Network
	network = &Network{
		events: make(map[string]func(clientId string, data []byte)),
	}
	return network
}

func (n *Network) StartServer(port int) {
	n.server = NewServer()
	n.server.Listen(port)
}

func (n *Network) ConnectClient(addr string) {
	n.client = NewClient()
	n.client.Connect(addr)
	n.client.WriteMessage("", []byte{})
}

func (n *Network) Update(dt float64) {
	if n.IsClient() {
		for packet, ok := n.client.GetNextMessage(); ok; packet, ok = n.client.GetNextMessage() {
			n.CallEvent(packet.Command, packet.Token, packet.Data)
		}
	} else if n.IsServer() {
		n.server.Update(dt)
		for packet, ok := n.server.GetNextMessage(); ok; packet, ok = n.server.GetNextMessage() {
			n.CallEvent(packet.Command, packet.Token, packet.Data)
		}
	}
}

func (n *Network) CallEvent(name, clientId string, data []byte) {
	callback, ok := n.events[name]
	if ok {
		log.Printf("[NETWORK EVENT] %v clientId:%v", name, clientId)
		callback(clientId, data)
	} else {
		log.Printf("[NETWORK EVENT] ERROR: Unknown event: %v clientId:%v", name, clientId)
	}
}

// RegisterEvent - register an event that will be triggered on clients and server.
func (n *Network) RegisterEvent(name string, callback func(clientId string, data []byte)) {
	n.events[name] = callback
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
	n.CallEvent(name, "", data)
	if n.IsClient() {
		n.client.WriteMessage(name, data)
	}
	if n.IsServer() {
		n.server.BroadcastMessage(name, data)
	}
}

//ClientJoinedEvent - This function is called when a new client connects to the server
func (n *Network) ClientJoinedEvent(callback func(clientId string)) {
	if n.server != nil {
		n.server.ClientJoinedEvent(callback)
	}
}

func (n *Network) ClientToken() string {
	return n.client.token
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
}

func (n *Network) KillClient() {
	if n.client != nil {
		n.client.Close()
	}
	n.client = nil
}

func (n *Network) KillServer() {
	if n.server != nil {
		n.server.Close()
	}
	n.server = nil
}
