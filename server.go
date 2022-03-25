package main

type Server struct {
	// clients hold all current clients in the chat server
	clients map[*Client]bool
	// join is a channel for clients wanting to join the server
	join chan *Client
	// leave is a channel for clients wanting to leave the server
	leave chan *Client
	// forward is a channel that listens for messages and forwards
	// them to other clients
	forward chan Message
}

// NewServer returns a new Server
func NewServer() *Server {
	return &Server{
		clients: make(map[*Client]bool),
		leave:   make(chan *Client),
		join:    make(chan *Client),
		forward: make(chan Message),
	}
}

func (s *Server) run() {
	for {
		select {
		case client := <-s.join:
			s.registerClient((client))
		case client := <-s.leave:
			s.unregisterClient(client)
		case msg := <-s.forward:
			s.broadcastToClients(msg)
		}
	}
}

func (s *Server) registerClient(c *Client) {
	s.clients[c] = true
}

func (s *Server) unregisterClient(c *Client) {
	if _, ok := s.clients[c]; ok {
		delete(s.clients, c)
	}
}

func (s *Server) broadcastToClients(msg Message) {
	for client := range s.clients {
		client.buffer <- msg
	}
}
