package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, err := strconv.ParseInt(q.Get("limit"), 10, 0)
	if err != nil {
		log.Print(err)
	}
	queryResult, err := find(q.Get("to"), q.Get("from"), limit)
	if err != nil {
		log.Print(err)
	}

	res, _ := json.Marshal(queryResult)
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func main() {
	wsServer := NewServer()
	go wsServer.run()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", HealthHandler)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(wsServer, w, r)
	})
	mux.HandleFunc("/messages", GetMessagesHandler)

	log.Printf("server listetening at localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
