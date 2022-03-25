package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	READ_BUFFER_SIZE  = 4096
	WRITE_BUFFER_SIZE = 1024
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  READ_BUFFER_SIZE,
	WriteBufferSize: WRITE_BUFFER_SIZE,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client represents an intermediary between a websocket connection
// and the server
type Client struct {
	conn   *websocket.Conn
	server *Server
	buffer chan Message
}

func newClient(conn *websocket.Conn, server *Server) *Client {
	return &Client{
		conn:   conn,
		server: server,
		buffer: make(chan Message),
	}
}

// ServeWs is a HTTP handler that upgrades the connection to the
// WebSocket protocol and creates a client.
func ServeWs(s *Server, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	client := newClient(conn, s)
	s.join <- client
	defer func() { s.leave <- client }()
	go client.write()
	go client.read()
}

// read transfers inbound messages from the websocket connection
// to the server
func (c *Client) read() {
	defer c.conn.Close()

	for {
		var msg Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			return
		}
		c.server.forward <- msg
	}

}

// write transfers messages from buffer channel to websocket
// connection
func (c *Client) write() {
	defer c.conn.Close()

	for msg := range c.buffer {
		err := c.conn.WriteJSON(msg)
		if err != nil {
			log.Print(err)
		}
	}
}
