package main

import (
	"ecommerce-backend/pkg/config"
	"ecommerce-backend/pkg/logger"
	"ecommerce-backend/services/orderservice/internal/handler"
	model "ecommerce-backend/services/orderservice/internal/models"
	"ecommerce-backend/services/orderservice/internal/repository"
	"ecommerce-backend/services/orderservice/internal/service"
	"log"

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

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	// AutoMigrate Order + OrderItem
	if err := db.AutoMigrate(&model.Order{}, &model.OrderItem{}); err != nil {
		log.Fatalf("❌ auto migrate failed: %v", err)
	}

	// Gin setup
	r := gin.Default()
	r.GET("/orderhealth", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "orderservice up"})
	})

	// after AutoMigrate
	orderRepo := repository.NewOrderRepository(db)
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
