package main

import (
	"log"
	"os"

	"ecommerce-backend/pkg/config"
	"ecommerce-backend/pkg/db"
	"ecommerce-backend/pkg/logger"
	"ecommerce-backend/pkg/middleware"
	"ecommerce-backend/services/cartservice/internal/handler"
	"ecommerce-backend/services/cartservice/internal/models"
	"ecommerce-backend/services/cartservice/internal/repository"
	"ecommerce-backend/services/cartservice/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// -----------------------------------------
	// 1. Init Logger
	// -----------------------------------------
	logger.Init()
	defer logger.Sync()

	// -----------------------------------------
	// 2. Load Environment Variables (.env optional)
	// -----------------------------------------
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("⚠️  No .env found, continuing with system environment variables")
	}

	// -----------------------------------------
	// 3. Database Setup
	// -----------------------------------------
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("❌ DATABASE_DSN environment variable not set")
	}

	dbConn, err := db.InitDB(dsn)
	if err != nil {
		log.Fatalf("❌ Failed to connect to DB: %v", err)
	}

	// Auto-migrate tables
	if err := dbConn.AutoMigrate(&models.Cart{}, &models.CartItem{}); err != nil {
		log.Fatalf("❌ AutoMigrate failed: %v", err)
	}

	// -----------------------------------------
	// 5. Dependency Injection
	// -----------------------------------------
	repo := repository.NewCartRepository(dbConn)
	cartService := service.NewCartService(repo, config.GetEnv("ORDER_SERVICE_URL", "http://localhost:8083"))
	cartHandler := handler.NewCartHandler(cartService)

	// -----------------------------------------
	// 6. Router Setup
	// -----------------------------------------
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "cartservice up"})
	})

	api := router.Group("/api/v1/cart")
	api.Use(middleware.JWTAuth())

	{
		api.POST("/items", cartHandler.AddItem)
		api.GET("", cartHandler.GetCart)
		api.PUT("/items/:id", cartHandler.UpdateItem)
		api.DELETE("/items/:id", cartHandler.DeleteItem)
		api.POST("/checkout", cartHandler.Checkout)
	}

	// -----------------------------------------
	// 7. Start Server
	// -----------------------------------------
	port := config.GetEnv("PORT", "8085")

	log.Printf("🚀 CartService running on port %s", port)
	router.Run(":" + port)
}
