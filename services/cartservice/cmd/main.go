package main

import (
	"ecommerce-backend/pkg/config"
	"ecommerce-backend/pkg/db"
	"ecommerce-backend/pkg/logger"
	"ecommerce-backend/services/cartservice/internal/handler"
	"ecommerce-backend/services/cartservice/internal/model"
	"ecommerce-backend/services/cartservice/internal/repository"
	"ecommerce-backend/services/cartservice/internal/service"
	"os"

	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	//config.LoadEnv()
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

	// AutoMigrate
	if err := gormDB.AutoMigrate(&model.Cart{}, &model.CartItem{}); err != nil {
		log.Fatalf("❌ auto migrate failed: %v", err)
	}

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "cartservice up"})
	})

	repo := repository.NewCartRepository(gormDB)
	svc := service.NewCartService(repo, "http://orderservice:8083") // ✅ internal Docker DNS
	handler := handler.NewCartHandler(svc)

	api := r.Group("/api/v1/cart")
	api.POST("/items", handler.AddItem)
	api.GET("", handler.GetCart)
	api.PUT("/items/:id", handler.UpdateItem)
	api.DELETE("/items/:id", handler.DeleteItem)
	api.POST("/checkout", handler.Checkout)

	port := config.GetEnv("PORT", "8085")
	log.Println("✅ CartService running on port", port)
	r.Run(":" + port)
}
