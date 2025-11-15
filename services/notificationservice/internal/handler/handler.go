package handler

import (
	"ecommerce-backend/services/notificationservice/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *service.NotificationService
}

func NewHandler(s *service.NotificationService) *Handler {
	return &Handler{svc: s}
}

type notifyReq struct {
	Type  string                 `json:"type" binding:"required"`
	To    string                 `json:"to" binding:"required,email"`
	Title string                 `json:"title" binding:"required"`
	Body  string                 `json:"body" binding:"required"`
	Meta  map[string]interface{} `json:"meta"`
}

func (h *Handler) Notify(c *gin.Context) {
	var req notifyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// accept request and enqueue for async processing
	if err := h.svc.EnqueueNotification(req.Type, req.To, req.Title, req.Body, req.Meta); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enqueue"})
		return
	}

	// 202 Accepted
	c.JSON(http.StatusAccepted, gin.H{"status": "accepted"})
}
