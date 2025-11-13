package service

import (
	model "ecommerce-backend/services/productservice/internal/models"
	repository "ecommerce-backend/services/productservice/internal/reposotory"
	"errors"
	"fmt"
	"time"
)

type ProductService struct {
	Repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{Repo: repo}
}

func (s *ProductService) CreateProduct(name, desc string, price float64, stock int) (*model.Product, error) {
	if name == "" || price <= 0 || stock < 0 {
		return nil, errors.New("invalid product details")
	}

	product := &model.Product{
		Name:        name,
		Description: desc,
		Price:       price,
		Stock:       stock,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.Repo.Create(product); err != nil {
		return nil, err
	}
	return product, nil
}

// GetAllProducts returns all products
func (s *ProductService) GetAllProducts() ([]model.Product, error) {
	return s.Repo.GetAll()
}

// GetProductByID returns a single product by ID
func (s *ProductService) GetProductByID(id string) (*model.Product, error) {
	return s.Repo.GetByID(id)
}

// Update an existing product
func (s *ProductService) UpdateProduct(id, name, desc string, price float64, stock int) (*model.Product, error) {
	product, err := s.Repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	if name != "" {
		product.Name = name
	}
	if desc != "" {
		product.Description = desc
	}
	if price > 0 {
		product.Price = price
	}
	if stock >= 0 {
		product.Stock = stock
	}

	if err := s.Repo.Update(product); err != nil {
		return nil, err
	}
	return product, nil
}

// Delete a product
func (s *ProductService) DeleteProduct(id string) error {
	return s.Repo.Delete(id)
}

func (s *ProductService) ReduceStock(productID string, qty int) error {
	product, err := s.Repo.GetByID(productID)
	if err != nil {
		return fmt.Errorf("failed to fetch product: %w", err)
	}
	if product == nil {
		return fmt.Errorf("product not found")
	}

	if product.Stock < qty {
		return fmt.Errorf("insufficient stock: available %d, requested %d", product.Stock, qty)
	}

	product.Stock -= qty

	if err := s.Repo.Update(product); err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	return nil
}
