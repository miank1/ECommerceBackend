package main

import (
	logger "ecommerce-backend/pkg/logger"
	"ecommerce-backend/services/userservice/internal/handler"
	"ecommerce-backend/services/userservice/internal/model"
	"ecommerce-backend/services/userservice/internal/repository"
	"ecommerce-backend/services/userservice/internal/service"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func main() {
	// Initialize global logger
	logger.Init()
	defer logger.Sync()

	// Load environment variables from .env file
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("⚠️  No .env file found, using system environment variables")
	}

	log.Println("Loaded DSN:", os.Getenv("DATABASE_DSN"))

	// Read environment variables
	dsn := os.Getenv("DATABASE_DSN")

	// Connect to Neon PostgreSQL with GORM
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			gormLogger.Config{
				SlowThreshold: time.Second,
				LogLevel:      gormLogger.Info, // change to Warn or Silent to reduce logs
				Colorful:      true,
			},
		),
	})
	if err != nil {
		log.Fatalf("❌ Failed to connect to PostgreSQL: %v\nDSN: %s", err, dsn)
	} else {
		log.Println("✅ Successfully connected to Neon PostgreSQL database.")
	}

	// Auto migrate User model
	if err := gormDB.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("❌ Auto migrate failed: %v", err)
	} else {
		log.Println("✅ User table migration successful!")
	}

	// Initialize repository, service, handler
	repo := repository.NewUserRepository(gormDB)
	svc := service.NewUserService(repo)
	h := handler.NewUserHandler(svc)

	// Initialize Gin router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "userservice is up"})
	})

	// Database connectivity check
	r.GET("/health/db", func(c *gin.Context) {
		sqlDB, err := gormDB.DB()
		if err != nil {
			c.JSON(500, gin.H{"db": "error", "details": err.Error()})
			return
		}
		if err := sqlDB.Ping(); err != nil {
			c.JSON(500, gin.H{"db": "not reachable", "details": err.Error()})
			return
		}
		c.JSON(200, gin.H{"db": "connected ✅"})
	})

	// User API routes
	api := r.Group("/api/v1")
	api.POST("/register", h.Register)
	api.POST("/login", h.Login)
	api.GET("/me", h.Me)

	// Start server
	log.Println("✅ UserService running on port 8081")
	if err := r.Run(":" + "8081"); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
