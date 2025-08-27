package main

import (
	"ecommerce-api/services/product/internal/handler"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/products", handler.GetProducts).Methods("GET")
	router.HandleFunc("/products", handler.CreateProduct).Methods("POST")

	// Start HTTP server
	log.Printf("Product service starting on port 8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
