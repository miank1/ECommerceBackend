package handler

import (
	"ecommerce-backend/services/orderservice/internal/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	Svc *service.OrderService
}

type updateOrderReq struct {
	Status string `json:"status" binding:"required"`
}

func NewOrderHandler(s *service.OrderService) *OrderHandler {
	return &OrderHandler{Svc: s}
}

type createOrderReq struct {
	Items  []service.OrderItemInput `json:"items" binding:"required"`
	UserID string                   `json:"user_id"`
}

func (h *OrderHandler) Create(c *gin.Context) {
	fmt.Println("Hello 2 ----------------------")
	var req createOrderReq
	// In real app, extract userID from JWT

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := req.UserID

	if userID == "" {
		// For now, fake a user until JWT is wired
		userID = "dummy-user-id"
	}
	order, err := h.Svc.CreateOrder(userID, req.Items)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"order": order})
}

// List orders for current user
func (h *OrderHandler) List(c *gin.Context) {
	// In real app, extract user_id from JWT
	userID := c.GetString("user_id")
	if userID == "" {
		userID = "dummy-user-id"
	}

	orders, err := h.Svc.GetOrdersByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

// Get a single order by ID
func (h *OrderHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing order id"})
		return
	}

	order, err := h.Svc.GetOrderByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch order"})
		return
	}
	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order})
}

// Update order status
func (h *OrderHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req updateOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.Svc.UpdateOrderStatus(id, req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order})
}

// Delete order
func (h *OrderHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.Svc.DeleteOrder(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order deleted"})
}
