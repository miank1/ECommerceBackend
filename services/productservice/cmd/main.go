package main

import (
	"ecommerce-backend/pkg/config"
	"ecommerce-backend/pkg/db"
	"ecommerce-backend/pkg/logger"
	"ecommerce-backend/services/productservice/internal/handler"
	model "ecommerce-backend/services/productservice/internal/models"
	repository "ecommerce-backend/services/productservice/internal/reposotory"
	"ecommerce-backend/services/productservice/internal/service"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Get database configuration from environment
	logger.Init()
	defer logger.Sync()

	// _ = godotenv.Load()

	cfg := db.Config{
		DSN:         os.Getenv("DATABASE_DSN"),
		MaxRetries:  6,
		RetryDelay:  2 * time.Second,
		ConnTimeout: 5 * time.Second,
	}

	gormDB, err := db.InitPostgres(cfg)
	if err != nil {
		log.Fatalf("could not initialize database: %v", err)
	}

	if err = gormDB.AutoMigrate(&model.Product{}); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	// Set up HTTP server
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "productservice up"})
	})

	// Wire dependencies (repo → service → handler)

	productRepo := repository.NewProductRepository(gormDB)
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
