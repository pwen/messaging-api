package main

import (
	"log"
	"net/http"
	"time"
)

// Message represents a message a client sends to the server
type Message struct {
	To        string    `json:"to"`
	From      string    `json:"from"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}

func main() {
	wsServer := NewServer()
	go wsServer.run()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", HealthHandler)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(wsServer, w, r)
	})

	log.Printf("server listetening at localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
