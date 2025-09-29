package main

import (
	"ecommerce-backend/pkg/config"
	"ecommerce-backend/pkg/db"
	"ecommerce-backend/pkg/logger"
	"ecommerce-backend/pkg/middleware"
	"ecommerce-backend/services/userservice/internal/handler"
	"ecommerce-backend/services/userservice/internal/model"
	"ecommerce-backend/services/userservice/internal/repository"
	"ecommerce-backend/services/userservice/internal/service"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	logger.Init()
	defer logger.Sync()

	// --- DB connection with retry ---

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

	// migrate User model
	if err := gormDB.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("❌ auto migrate failed: %v", err)
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
		c.JSON(200, gin.H{"status": "userservice up"})
	})

	api := r.Group("/api/v1")
	api.POST("/register", h.Register)
	api.POST("/login", h.Login)

	protected := api.Group("")
	protected.Use(middleware.JWTAuth())
	protected.GET("/me", h.Me)

	port := config.GetEnv("PORT", "8081")
	log.Println("✅ UserService running on port", port)
	r.Run(":" + port)
}
