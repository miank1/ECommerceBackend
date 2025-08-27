package repository

import (
	"context"
	"ecommerce-api/services/product/internal/models"
)

type ProductRepository interface {
	GetProducts(ctx context.Context) ([]models.Product, error)
	GetProductByID(ctx context.Context, id string) (*models.Product, error)
	CreateProduct(ctx context.Context, product *models.Product) error
	UpdateProduct(ctx context.Context, product *models.Product) error
	DeleteProduct(ctx context.Context, id string) error
}
