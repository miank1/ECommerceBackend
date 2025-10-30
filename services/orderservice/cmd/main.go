package main

import (
	"ecommerce-backend/pkg/config"
	"ecommerce-backend/pkg/db"
	"ecommerce-backend/pkg/logger"
	"ecommerce-backend/services/orderservice/internal/handler"
	model "ecommerce-backend/services/orderservice/internal/models"
	"ecommerce-backend/services/orderservice/internal/repository"
	"ecommerce-backend/services/orderservice/internal/service"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Initialize global logger

	logger.Init()
	defer logger.Sync()

	// Load environment variables from .env file
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("⚠️  No .env file found, using system environment variables")
	}

	dsn := os.Getenv("DATABASE_DSN")

	gormDB, err := db.InitDB(dsn)
	if err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}

	// AutoMigrate Order + OrderItem
	if err := gormDB.AutoMigrate(&model.Order{}, &model.OrderItem{}); err != nil {
		log.Fatalf("❌ auto migrate failed: %v", err)
	}

	// Gin setup
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "orderservice up"})
	})

	// after AutoMigrate
	orderRepo := repository.NewOrderRepository(gormDB)
	orderSvc := service.NewOrderService(orderRepo)
	orderHandler := handler.NewOrderHandler(orderSvc)

	// routes

	api := r.Group("/api/v1")
	api.POST("/orders", orderHandler.Create)
	api.GET("/orders", orderHandler.List)
	api.GET("/orders/:id", orderHandler.GetByID)
	api.PUT("/orders/:id", orderHandler.Update)
	api.DELETE("/orders/:id", orderHandler.Delete)

	port := config.GetEnv("PORT", "8083")
	log.Println("✅ OrderService running on port", port)
	r.Run(":" + port)
}
