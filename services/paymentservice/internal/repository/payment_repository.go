package repository

import (
	"ecommerce-backend/services/paymentservice/internal/models"

	"gorm.io/gorm"
)

type PaymentRepository struct {
	DB *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{DB: db}
}

func (r *PaymentRepository) Create(payment *models.Payment) error {
	return r.DB.Create(payment).Error
}

func (r *PaymentRepository) GetByOrderID(orderID string) (*models.Payment, error) {
	var payment models.Payment
	if err := r.DB.Where("order_id = ?", orderID).First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}
