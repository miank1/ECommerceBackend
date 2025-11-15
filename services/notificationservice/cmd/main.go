package main

import (
	"ecommerce-backend/services/notificationservice/internal/handler"
	"ecommerce-backend/services/notificationservice/internal/service"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// config from env
	smtpHost := os.Getenv("SMTP_HOST")  // e.g., "mailhog"
	smtpPort := os.Getenv("SMTP_PORT")  // e.g., "1025"
	smtpFrom := os.Getenv("NOTIF_FROM") // e.g., "no-reply@local.test"
	workerCount := 3

	notifSvc := service.NewNotificationService(smtpHost+":"+smtpPort, smtpFrom)
	// Start background worker(s)
	notifSvc.StartWorkers(workerCount)

	r := gin.Default()
	h := handler.NewHandler(notifSvc)
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	api := r.Group("/api/v1")
	{
		api.POST("/notify", h.Notify) // accepts notification requests
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}
	log.Println("Notification service running on port", port)
	r.Run(":" + port)
}
