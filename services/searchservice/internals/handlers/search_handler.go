package handler

import (
	"ecommerce-backend/services/searchservice/internals/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	service service.SearchService
}

func NewSearchHandler(s service.SearchService) *SearchHandler {
	return &SearchHandler{service: s}
}

func (h *SearchHandler) Search(c *gin.Context) {
	q := c.Query("q")
	category := c.Query("category")
	minPriceStr := c.Query("minPrice")
	maxPriceStr := c.Query("maxPrice")

	var minPrice, maxPrice float64
	var err error

	if minPriceStr != "" {
		minPrice, err = strconv.ParseFloat(minPriceStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid minPrice"})
			return
		}
	}
	if maxPriceStr != "" {
		maxPrice, err = strconv.ParseFloat(maxPriceStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid maxPrice"})
			return
		}
	}

	products, err := h.service.SearchProducts(q, category, minPrice, maxPrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search products"})
		return
	}

	c.JSON(http.StatusOK, products)
}
