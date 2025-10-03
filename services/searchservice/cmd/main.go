package main

import (
	handler "ecommerce-backend/services/searchservice/internals/handlers"
	"ecommerce-backend/services/searchservice/internals/repository"
	"ecommerce-backend/services/searchservice/internals/service"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("DATABASE_DSN")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	// Connect DB
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to ProductDB: %v", err)
	}

	// Setup layers
	repo := repository.NewProductRepository(db)
	svc := service.NewSearchService(repo)
	handler := handler.NewSearchHandler(svc)

	http.HandleFunc("/search", handler.Search)

	log.Printf("✅ Search Service running on :%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
