package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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
var (
	products = make(map[string]Product)
	nextID   = 1 // For auto-incrementing IDs
)

// GetProducts returns all products
func GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Convert map to slice for JSON response
	productList := make([]Product, 0, len(products))
	for _, product := range products {
		productList = append(productList, product)
	}

	for x, v := range productList {
		fmt.Printf("Product %d: %+v\n", x, v)
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

	// Set auto-incrementing ID
	product.ID = fmt.Sprintf("%d", nextID)
	products[product.ID] = product
	nextID++ // Increment for next product ID

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("Getting single product")
	// Get product ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Find product
	product, exists := products[id]
	fmt.Println("Product found:", product)
	if !exists {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(product)
}

// UpdateProduct updates an existing product
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get product ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if product exists
	_, exists := products[id]
	if !exists {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Parse updated product from request body
	var updatedProduct Product
	if err := json.NewDecoder(r.Body).Decode(&updatedProduct); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate updated product
	if updatedProduct.Name == "" {
		http.Error(w, "Product name is required", http.StatusBadRequest)
		return
	}

	if updatedProduct.Price <= 0 {
		http.Error(w, "Price must be greater than 0", http.StatusBadRequest)
		return
	}

	// Keep the same ID, update other fields
	updatedProduct.ID = id
	products[id] = updatedProduct

	json.NewEncoder(w).Encode(updatedProduct)
}

// DeleteProduct removes a product
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get product ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if product exists
	_, exists := products[id]
	if !exists {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Delete the product
	delete(products, id)

	// Return success message
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Product deleted successfully"})
}
