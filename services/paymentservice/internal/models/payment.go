package models

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	OrderID       uuid.UUID `gorm:"not null" json:"order_id"`
	UserID        uuid.UUID `gorm:"not null" json:"user_id"`
	Amount        float64   `gorm:"not null" json:"amount"`
	Status        string    `gorm:"not null" json:"status"` // pending, success, failed
	PaymentMethod string    `gorm:"not null" json:"payment_method"`
	TransactionID string    `gorm:"uniqueIndex" json:"transaction_id"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
