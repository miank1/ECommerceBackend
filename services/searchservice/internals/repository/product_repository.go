package repository

import (
	"ecommerce-backend/services/searchservice/internals/models"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Search(q, category string, minPrice, maxPrice float64) ([]models.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Search(q, category string, minPrice, maxPrice float64) ([]models.Product, error) {
	var products []models.Product
	query := r.db

	if q != "" {
		query = query.Where("LOWER(name) LIKE ?", "%"+q+"%")
	}
	if category != "" {
		query = query.Where("LOWER(category) = LOWER(?)", category)
	}
	if minPrice > 0 {
		query = query.Where("price >= ?", minPrice)
	}
	if maxPrice > 0 {
		query = query.Where("price <= ?", maxPrice)
	}

	err := query.Find(&products).Error
	return products, err
}
