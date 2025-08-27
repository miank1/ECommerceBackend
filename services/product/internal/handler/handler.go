package handler

import (
	"encoding/json"
	"net/http"
)

// Product represents a single product in our store
type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

// Simple in-memory storage
var products = make(map[string]Product)

// GetProducts returns all products
func GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Convert map to slice for JSON response
	productList := make([]Product, 0, len(products))
	for _, product := range products {
		productList = append(productList, product)
	}

	json.NewEncoder(w).Encode(productList)
}

// createProduct adds a new product
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Simple validation
	if product.Name == "" {
		http.Error(w, "Product name is required", http.StatusBadRequest)
		return
	}

	if product.Price <= 0 {
		http.Error(w, "Price must be greater than 0", http.StatusBadRequest)
		return
	}

	// Store the product (using Name as ID for now)
	product.ID = product.Name // In real app, use UUID
	products[product.ID] = product

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}
