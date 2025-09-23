package main

import (
	"ecommerce-backend/pkg/config"
	"ecommerce-backend/pkg/logger"
	"ecommerce-backend/services/orderservice/internal/handler"
	model "ecommerce-backend/services/orderservice/internal/models"
	"ecommerce-backend/services/orderservice/internal/repository"
	"ecommerce-backend/services/orderservice/internal/service"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	//config.LoadEnv()
	logger.Init()
	defer logger.Sync()

	dsn := config.GetEnv("DATABASE_DSN", "")
	log.Println("Using DSN:", dsn)

	// Retry DB connect
	var db *gorm.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("✅ Connected to DB Order Service ")
			break
		}
		log.Println("DB not ready, retrying in 3s...", err)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		log.Fatalf("❌ failed to connect db after retries: %v", err)
	}

	// AutoMigrate Order + OrderItem
	if err := db.AutoMigrate(&model.Order{}, &model.OrderItem{}); err != nil {
		log.Fatalf("❌ auto migrate failed: %v", err)
	}

	// Gin setup
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "orderservice up"})
	})

	// after AutoMigrate
	orderRepo := repository.NewOrderRepository(db)
	orderSvc := service.NewOrderService(orderRepo)
	orderHandler := handler.NewOrderHandler(orderSvc)

	// TODO: add order routes
	// routes

	api := r.Group("/api/v1")
	api.POST("/orders", orderHandler.Create)
	api.GET("/orders", orderHandler.List)
	api.GET("/orders/:id", orderHandler.GetByID)
	api.POST("/orders", orderHandler.Create)

	port := config.GetEnv("PORT", "8083")
	log.Println("✅ OrderService running on port", port)
	r.Run(":" + port)
}
