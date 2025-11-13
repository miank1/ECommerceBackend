package main

import (
	"ecommerce-backend/pkg/config"
	"ecommerce-backend/pkg/db"
	"ecommerce-backend/pkg/logger"
	"ecommerce-backend/services/productservice/internal/handler"
	model "ecommerce-backend/services/productservice/internal/models"
	repository "ecommerce-backend/services/productservice/internal/reposotory"
	"ecommerce-backend/services/productservice/internal/service"
	"ecommerce-backend/services/productservice/seed"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Get database configuration from environment
	logger.Init()
	defer logger.Sync()

	// Load environment variables from .env file
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("⚠️  No .env file found, using system environment variables")
	}

	log.Println("Loaded DSN:", os.Getenv("DATABASE_DSN"))

	dsn := os.Getenv("DATABASE_DSN")

	gormDB, err := db.InitDB(dsn)
	if err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}

	if err = gormDB.AutoMigrate(&model.Product{}); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	var count int64

	gormDB.Model(&repository.ProductRepository{}).Count(&count)

	if count == 0 {
		log.Println("Seeding products for the first time ... ")
		seed.SeedProducts(gormDB)
		log.Println("Products seeded successfully ")
	} else {
		log.Println("Products alreadys existed, skipping seed.")
	}

	// Set up HTTP server
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "productservice up"})
	})

	// Repository Pattern

	productRepo := repository.NewProductRepository(gormDB)
	productSvc := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productSvc)

	api := r.Group("/api/v1")
	api.POST("/products", productHandler.Create)
	api.GET("/products", productHandler.List)
	api.GET("/products/:id", productHandler.GetByID)
	api.PUT("/products/:id", productHandler.Update)
	api.DELETE("/products/:id", productHandler.Delete)
	api.PATCH("/products/:id/reduce-stock", productHandler.ReduceStock)

	port := config.GetEnv("PORT", "8082")
	fmt.Println("✅ ProductService running on port", port)
	r.Run(":" + port)
}
