package main

import (
	"ecommerce-backend/pkg/db"
	"ecommerce-backend/pkg/logger"
	"ecommerce-backend/services/paymentservice/internal/handler"
	"ecommerce-backend/services/paymentservice/internal/models"
	"ecommerce-backend/services/paymentservice/internal/repository"
	"ecommerce-backend/services/paymentservice/internal/service"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
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

	// Auto migrate Payment model
	gormDB.AutoMigrate(&models.Payment{})

	// Initialize layers
	repo := repository.NewPaymentRepository(gormDB)
	svc := service.NewPaymentService(repo)
	handler := handler.NewPaymentHandler(svc)

	// Setup routes

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "payment service is service up"})
	})

	api := r.Group("/api/v1")
	{
		api.POST("/payments", handler.ProcessPayment)
	}

	log.Println("🚀 Payment Service running on :8086 ***********")
	r.Run(":8086")
}
