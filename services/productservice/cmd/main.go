package main

import (
	"ecommerce-backend/pkg/config"
	"ecommerce-backend/pkg/logger"
	"ecommerce-backend/services/productservice/internal/handler"
	model "ecommerce-backend/services/productservice/internal/models"
	repository "ecommerce-backend/services/productservice/internal/reposotory"
	"ecommerce-backend/services/productservice/internal/service"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Get database configuration from environment
	logger.Init()
	defer logger.Sync()

	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("DATABASE_DSN environment variable is required")
	}

	log.Printf("Connecting to database with DSN: %s\n", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	// Auto migrate Product model
	if err = db.AutoMigrate(&model.Product{}); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	// Set up HTTP server
	r := gin.Default()

	// Health check endpoint
	r.GET("/producthealth", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "productservice up"})
	})

	// Wire dependencies (repo → service → handler)
	// TODO: add ProductService + ProductHandler when we build APIs

	productRepo := repository.NewProductRepository(db)
	productSvc := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productSvc)

	api := r.Group("/api/v1")
	api.POST("/products", productHandler.Create)
	api.GET("/products", productHandler.List)
	api.GET("/products/:id", productHandler.GetByID)
	api.PUT("/products/:id", productHandler.Update)
	api.DELETE("/products/:id", productHandler.Delete)

	port := config.GetEnv("PORT", "8082")
	log.Println("✅ ProductService running on port", port)
	r.Run(":" + port)
}
