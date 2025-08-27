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

// Validation constants
const (
	minNameLength        = 3
	maxNameLength        = 100
	maxDescriptionLength = 500
	minPrice             = 0.01
	maxPrice             = 999999.99
	maxStock             = 999999
)

// validateProduct performs all validation checks on a product
func validateProduct(p Product) error {
	// Name validation
	if p.Name == "" {
		return fmt.Errorf("product name is required")
	}
	if len(p.Name) < minNameLength {
		return fmt.Errorf("product name must be at least %d characters", minNameLength)
	}
	if len(p.Name) > maxNameLength {
		return fmt.Errorf("product name cannot exceed %d characters", maxNameLength)
	}

	// Description validation
	if len(p.Description) > maxDescriptionLength {
		return fmt.Errorf("description cannot exceed %d characters", maxDescriptionLength)
	}

	// Price validation
	if p.Price < minPrice {
		return fmt.Errorf("price must be at least %.2f", minPrice)
	}
	if p.Price > maxPrice {
		return fmt.Errorf("price cannot exceed %.2f", maxPrice)
	}

	// Stock validation
	if p.Stock < 0 {
		return fmt.Errorf("stock cannot be negative")
	}
	if p.Stock > maxStock {
		return fmt.Errorf("stock cannot exceed %d", maxStock)
	}

	return nil
}

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

	// Comprehensive validation
	if err := validateProduct(product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check for unique name
	for _, p := range products {
		if p.Name == product.Name {
			http.Error(w, "Product with this name already exists", http.StatusConflict)
			return
		}
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
	if err := validateProduct(updatedProduct); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check for unique name, excluding the current product
	for pid, p := range products {
		if p.Name == updatedProduct.Name && pid != id {
			http.Error(w, "Product with this name already exists", http.StatusConflict)
			return
		}
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
