package main

import (
	"ecommerce-api/services/user/handlers"
	"ecommerce-api/services/user/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// Public routes
	router.HandleFunc("/users/register", handlers.RegisterUser).Methods("POST")
	router.HandleFunc("/users/login", handlers.LoginUser).Methods("POST")

	// Protected routes (with JWT authentication)
	router.HandleFunc("/users/profile", middleware.AuthMiddleware(handlers.GetProfile)).Methods("GET")

	// Start HTTP server
	log.Printf("User service starting on port 8081...")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
