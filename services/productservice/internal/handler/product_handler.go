package handler

import (
	"ecommerce-backend/services/productservice/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	Svc *service.ProductService
}

// ReduceStockRequest defines the request structure
type ReduceStockRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required"`
}

func NewProductHandler(s *service.ProductService) *ProductHandler {
	return &ProductHandler{Svc: s}
}

type createProductReq struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
	Stock       int     `json:"stock" binding:"required"`
}

type updateProductReq struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

func (h *ProductHandler) Create(c *gin.Context) {

	var req createProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	product, err := h.Svc.CreateProduct(req.Name, req.Description, req.Price, req.Stock)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "product": product})
}

// List all products
func (h *ProductHandler) List(c *gin.Context) {
	products, err := h.Svc.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to fetch products"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "products": products})
}

// Get product by ID
func (h *ProductHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "missing product id"})
		return
	}

	product, err := h.Svc.GetProductByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to fetch product"})
		return
	}
	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "product": product})
}

// Update product
func (h *ProductHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req updateProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	product, err := h.Svc.UpdateProduct(id, req.Name, req.Description, req.Price, req.Stock)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "product": product})
}

// Delete product
func (h *ProductHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.Svc.DeleteProduct(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "product deleted"})
}

func (h *ProductHandler) ReduceStock(c *gin.Context) {
	var req ReduceStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Svc.ReduceStock(req.ProductID, req.Quantity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Stock reduced successfully",
		"product_id":    req.ProductID,
		"quantity_sold": req.Quantity,
	})
}
