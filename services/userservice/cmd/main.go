package main

import (
	"ecommerce-backend/pkg/config"
	logger "ecommerce-backend/pkg/logger"
	"ecommerce-backend/pkg/middleware"
	"ecommerce-backend/services/userservice/internal/handler"
	"ecommerce-backend/services/userservice/internal/model"
	"ecommerce-backend/services/userservice/internal/repository"
	"ecommerce-backend/services/userservice/internal/service"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	logger.Init()
	defer logger.Sync()

	dsn := os.Getenv("DATABASE_DSN")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	// Connect DB
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to ProductDB: %v", err)
	}
	// migrate User model
	if err := gormDB.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("❌ auto migrate failed: %v", err)
	} else {
		log.Println(" ✅ Migration Successful user service!!")
	}

	// wire layers
	repo := repository.NewUserRepository(gormDB)
	ttlStr := config.GetEnv("TOKEN_TTL_MIN", "60")
	ttl, _ := strconv.Atoi(ttlStr)
	secret := config.GetEnv("JWT_SECRET", "changeme")
	svc := service.NewUserService(repo, secret, ttl)
	h := handler.NewUserHandler(svc)

	// gin setup
	r := gin.Default()

	// health
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "userservice is up"})
	})

	api := r.Group("/api/v1")
	api.POST("/register", h.Register)
	api.POST("/login", h.Login)

	protected := api.Group("")
	protected.Use(middleware.JWTAuth())
	protected.GET("/me", h.Me)

	port = config.GetEnv("PORT", "8081")
	log.Println("✅ UserService running on port", port)
	r.Run(":" + port)
}
