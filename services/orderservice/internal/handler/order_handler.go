package handler

import (
	"ecommerce-backend/services/orderservice/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	Svc *service.OrderService
}

func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{Svc: svc}
}

type CreateOrderRequest struct {
	UserID string                 `json:"user_id"`
	Items  []service.OrderItemReq `json:"items"`
}

// POST /api/v1/orders
func (h *OrderHandler) Create(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.UserID == "" || len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order payload"})
		return
	}

	order, err := h.Svc.CreateOrder(req.UserID, req.Items)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"order": order})
}

// GET /api/v1/orders/:id
func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order id required"})
		return
	}

	order, err := h.Svc.GetOrderByID(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order})
}

// PATCH /api/v1/orders/:id/status
func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	orderID := c.Param("id")

	var req struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	order, err := h.Svc.UpdateStatus(orderID, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("✅ Order status updated: ****************", orderID, "to", req.Status)

	// ✅ Trigger inventory update when payment succeeded
	if req.Status == "paid" {
		go func() {
			if err := h.Svc.UpdateInventory(orderID); err != nil {
				log.Println("⚠️ Inventory update failed:", err)
			} else {
				log.Println("✅ Inventory updated for order:", orderID)
			}
		}()
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Order status updated successfully",
		"order_id":   orderID,
		"new_status": order.Status,
	})
}
