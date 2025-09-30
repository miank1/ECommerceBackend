package service

import (
	"ecommerce-backend/services/searchservice/internals/models"
	"ecommerce-backend/services/searchservice/internals/repository"
)

type SearchService interface {
	SearchProducts(q, category string, minPrice, maxPrice float64) ([]models.Product, error)
}

type searchService struct {
	repo repository.ProductRepository
}

func NewSearchService(repo repository.ProductRepository) SearchService {
	return &searchService{repo: repo}
}

func (s *searchService) SearchProducts(q, category string, minPrice, maxPrice float64) ([]models.Product, error) {
	return s.repo.Search(q, category, minPrice, maxPrice)
}
