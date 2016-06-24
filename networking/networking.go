package networking

import "log"

type Network struct {
	client *Client
	server *Server
	events map[string]func(clientId string, args ...interface{})
}

func NewNetwork() *Network {
	var network *Network
	network = &Network{
		events: make(map[string]func(clientId string, args ...interface{})),
	}
	return network
}

func (n *Network) StartServer(port int) {
	n.KillClient()
	n.server = NewServer()
	n.server.Listen(port)
}

func (n *Network) ConnectClient(addr string) {
	n.KillServer()
	n.client = NewClient()
	n.client.Connect(addr)
	n.client.WriteMessage("")
}

func (n *Network) Update(dt float64) {
	if n.IsClient() {
		// for packet, ok := n.client.GetNextMessage(); ok; {
		if packet, ok := n.client.GetNextMessage(); ok {
			n.CallEvent(packet.Command, packet.Token, packet.Args...)
		}
	} else if n.IsServer() {
		n.server.Update(dt)
		// for packet, ok := n.server.GetNextMessage(); ok; {
		if packet, ok := n.server.GetNextMessage(); ok {
			n.CallEvent(packet.Command, packet.Token, packet.Args...)
		}
	}
}

func (n *Network) CallEvent(name, clientId string, args ...interface{}) {
	callback, ok := n.events[name]
	if ok {
		log.Printf("[NETWORK EVENT] %v clientId:%v", name, clientId)
		callback(clientId, args...)
	}
}

// RegisterEvent - register an event that will be triggered on clients and server.
func (n *Network) RegisterEvent(name string, callback func(clientId string, args ...interface{})) {
	n.events[name] = callback
}

// TriggerEvent - Trigger an event to run on a particular client.
// If called on the client, this will trigger the event on the server.
func (n *Network) TriggerEvent(name, clientId string, args ...interface{}) {
	if n.IsClient() {
		n.client.WriteMessage(name, args...)
	}
	if n.IsServer() {
		n.server.WriteMessage(name, clientId, args...)
	}
}

// BroadcastEvent - trigger an event on all clients.
// If called on the client, this will trigger the event on the server.
func (n *Network) BroadcastEvent(name string, args ...interface{}) {
	if n.IsClient() {
		n.client.WriteMessage(name, args...)
	}
	if n.IsServer() {
		n.server.BroadcastMessage(name, args...)
	}
}

// CallOnServerAndClient - trigger an event on the server and on all client.
// If called on the client, this will trigger the event on the client and on the server.
func (n *Network) TriggerOnServerAndClients(name string, args ...interface{}) {
	n.CallEvent(name, "", args...)
	if n.IsClient() {
		n.client.WriteMessage(name, args...)
	}
	if n.IsServer() {
		n.server.BroadcastMessage(name, args...)
	}
}

//ClientJoinedEvent - This function is called when a new client connects to the server
func (n *Network) ClientJoinedEvent(callback func(clientId string)) {
	if n.server != nil {
		n.server.ClientJoinedEvent(callback)
	}
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
