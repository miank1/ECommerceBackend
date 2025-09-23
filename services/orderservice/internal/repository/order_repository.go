package repository

import (
	model "ecommerce-backend/services/orderservice/internal/models"

	"gorm.io/gorm"
)

type OrderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (r *OrderRepository) Create(order *model.Order) error {
	return r.DB.Create(order).Error
}

func (r *OrderRepository) GetByID(id string) (*model.Order, error) {
	var order model.Order
	if err := r.DB.Preload("Items").First(&order, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// List all orders for a user
func (r *OrderRepository) GetByUser(userID string) ([]model.Order, error) {
	var orders []model.Order
	if err := r.DB.Preload("Items").Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
