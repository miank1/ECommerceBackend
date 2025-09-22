package main

import (
	"ecommerce-backend/pkg/config"
	"ecommerce-backend/pkg/logger"
	"ecommerce-backend/pkg/middleware"
	"ecommerce-backend/services/userservice/internal/handler"
	"ecommerce-backend/services/userservice/internal/model"
	"ecommerce-backend/services/userservice/internal/repository"
	"ecommerce-backend/services/userservice/internal/service"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// load env + logger
	//config.LoadEnv()
	logger.Init()
	defer logger.Sync()

	dsn := config.GetEnv("DATABASE_DSN", "")
	log.Println("Using DSN:", dsn)

	// --- DB connection with retry ---
	var db *gorm.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}) // use "=" not ":="
		fmt.Println("Hello World !!!!!!!!!")
		if err == nil {
			log.Println("✅ Connected to DB")
			break
		}
		log.Println("DB not ready, retrying in 3s...", err)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		log.Fatalf("❌ failed to connect db after retries: %v", err)
	}
	if db == nil {
		log.Fatal("❌ db is nil after retries")
	}

	// migrate User model
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("❌ auto migrate failed: %v", err)
	}

	// wire layers
	repo := repository.NewUserRepository(db)
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
