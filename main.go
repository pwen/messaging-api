package main

import (
	"log"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", HealthHandler)

	log.Printf("server listetening at localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
