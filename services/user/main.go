package main

import (
	"ecommerce-api/services/user/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/users/register", handlers.RegisterUser).Methods("POST")
	router.HandleFunc("/users/login", handlers.LoginUser).Methods("POST")

	// Start HTTP server
	log.Printf("User service starting on port 8081...")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
