package repository

import (
	"ecommerce-backend/services/orderservice/internal/models"

	"gorm.io/gorm"
)

type OrderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (r *OrderRepository) Create(order *models.Order) error {
	return r.DB.Create(order).Error
}

func (r *OrderRepository) GetByID(id string) (*models.Order, error) {
	var order models.Order
	if err := r.DB.Preload("Items").First(&order, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) Save(order *models.Order) error {
	return r.DB.Save(order).Error
}
