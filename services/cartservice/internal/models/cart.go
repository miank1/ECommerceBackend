package models

import (
	"ecommerce-backend/services/searchservice/internals/models"
	"time"

	uuid "github.com/google/uuid"
)

// 🛒 Cart model — represents one user’s cart
type Cart struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    uuid.UUID  `gorm:"not null" json:"user_id"`
	Items     []CartItem `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE;" json:"items"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// CartItem — each product in the cart
type CartItem struct {
	ID        uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CartID    uuid.UUID       `gorm:"not null" json:"cart_id"`
	ProductID uuid.UUID       `gorm:"not null" json:"product_id"`
	Quantity  int             `gorm:"not null" json:"quantity"`
	Price     float64         `gorm:"default:0" json:"price"` // ✅ New field: store product price at checkout time
	Product   *models.Product `json:"product"`
	CreatedAt time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

type CheckoutEvent struct {
	CheckoutID string     `json:"checkout_id"`
	UserID     string     `json:"user_id"`
	Items      []CartItem `json:"items"`
	TotalPrice float64    `json:"total_price"`
	CreatedAt  time.Time  `json:"created_at"`
}
