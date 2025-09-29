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
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	logger.Init()
	defer logger.Sync()

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

	// AutoMigrate Order + OrderItem
	if err := gormDB.AutoMigrate(&model.Order{}, &model.OrderItem{}); err != nil {
		log.Fatalf("❌ auto migrate failed: %v", err)
	}

	// Gin setup
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "orderservice up"})
	})

	r.GET("/health1", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Hello World !!"})
	})

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Hello World !!!!"})
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
