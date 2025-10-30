package main

import (
	"ecommerce-backend/pkg/config"
	"ecommerce-backend/pkg/db"
	"ecommerce-backend/pkg/logger"
	"ecommerce-backend/pkg/middleware"
	"ecommerce-backend/services/cartservice/internal/handler"
	"ecommerce-backend/services/cartservice/internal/model"
	"ecommerce-backend/services/cartservice/internal/repository"
	"ecommerce-backend/services/cartservice/internal/service"
	"os"

	"log"

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

	// AutoMigrate
	if err := gormDB.AutoMigrate(&model.Cart{}, &model.CartItem{}); err != nil {
		log.Fatalf("❌ auto migrate failed: %v", err)
	}

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "cartservice up"})
	})

	repo := repository.NewCartRepository(gormDB)
	svc := service.NewCartService(repo, "http://localhost:8083") // ✅ internal Docker DNS
	handler := handler.NewCartHandler(svc)

	api := r.Group("/api/v1/cart")

	api.Use(middleware.JWTAuth())
	{
		api.POST("/items", handler.AddItem)
		api.GET("", handler.GetCart)
		api.PUT("/items/:id", handler.UpdateItem)
		api.DELETE("/items/:id", handler.DeleteItem)
		api.POST("/checkout", handler.Checkout)
	}

	port := config.GetEnv("PORT", "8085")
	log.Println("✅ CartService running on port -----------", port)
	r.Run(":" + port)
}
