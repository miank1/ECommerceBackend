package handler

import (
	"ecommerce-backend/services/productservice/internal/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	Svc *service.ProductService
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

func (h *ProductHandler) Create(c *gin.Context) {

	fmt.Println(" ✅ Hello World Product Service ✅")
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
