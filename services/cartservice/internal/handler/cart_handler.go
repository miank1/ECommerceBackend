package handler

import (
	"net/http"

	"ecommerce-backend/services/cartservice/internal/service"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	Svc *service.CartService
}

func NewCartHandler(s *service.CartService) *CartHandler {
	return &CartHandler{Svc: s}
}

type addItemReq struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required"`
}

type updateItemReq struct {
	Quantity int `json:"quantity" binding:"required"`
}

func (h *CartHandler) AddItem(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var req addItemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.Svc.AddItem(userID, req.ProductID, req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item added to cart", "cart": cart})
}

func (h *CartHandler) GetCart(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	cart, err := h.Svc.GetCart(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if cart == nil {
		c.JSON(http.StatusOK, gin.H{"cart": gin.H{"id": "", "items": []interface{}{}}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cart": cart})
}

func (h *CartHandler) UpdateItem(c *gin.Context) {
	itemID := c.Param("id")

	var req updateItemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.Svc.UpdateItemQuantity(itemID, req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item updated", "item": item})
}

func (h *CartHandler) DeleteItem(c *gin.Context) {
	itemID := c.Param("id")

	if err := h.Svc.DeleteItem(itemID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item removed from cart"})
}

func (h *CartHandler) Checkout(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	orderResp, err := h.Svc.Checkout(c, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "checkout successful",
		"order":   orderResp,
	})
}

func currentUserID(c *gin.Context) (string, bool) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authenticated user"})
		return "", false
	}

	return userID, true
}
