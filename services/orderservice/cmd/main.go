package main

import (
	"ecommerce-backend/pkg/db"
	"ecommerce-backend/pkg/logger"
	"ecommerce-backend/services/orderservice/internal/handler"
	"ecommerce-backend/services/orderservice/internal/models"
	"ecommerce-backend/services/orderservice/internal/repository"
	"ecommerce-backend/services/orderservice/internal/service"
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
	dsn := os.Getenv("DATABASE_DSN")

	gormDB, err := db.InitDB(dsn)
	if err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}

	if err := gormDB.AutoMigrate(&models.Order{}, &models.OrderItem{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	repo := repository.NewOrderRepository(gormDB)
	svc := service.NewOrderService(repo)
	h := handler.NewOrderHandler(svc)

	r := gin.Default()

	api := r.Group("/api/v1/orders")
	{
		api.POST("", h.Create)
		api.GET("/:id", h.GetOrder)
		api.POST("/:id/update-status", h.UpdateStatus)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}
	log.Printf("✅ OrderService running on port %s", port)
	r.Run(":" + port)
}
