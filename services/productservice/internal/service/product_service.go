package service

import (
	model "ecommerce-backend/services/productservice/internal/models"
	repository "ecommerce-backend/services/productservice/internal/reposotory"
	"errors"
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
