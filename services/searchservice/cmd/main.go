package main

import (
	"ecommerce-backend/pkg/config"
	handler "ecommerce-backend/services/searchservice/internals/handlers"
	"ecommerce-backend/services/searchservice/internals/repository"
	"ecommerce-backend/services/searchservice/internals/service"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("DATABASE_DSN env not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("cannot connect to DB: %v", err)
	}

	// Set up HTTP server
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "searchservice is up"})
	})

	repo := repository.NewProductRepository(db)
	svc := service.NewSearchService(repo)
	searchHandler := handler.NewSearchHandler(svc)

	api := r.Group("/api/v1")
	api.GET("/search", searchHandler.Search)

	port := config.GetEnv("PORT", "8084")
	log.Println("✅ SearchService running on port", port)
	r.Run(":" + port)
}
