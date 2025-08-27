package main

import (
	"log"
	"net/http"
)

func main() {
	// Initialize server
	port := ":8080"
	log.Printf("Starting server on port %s", port)

	// Create a new router
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start the server
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
