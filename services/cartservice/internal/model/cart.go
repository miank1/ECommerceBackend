package model

import "time"

// 🛒 Cart model — represents one user’s cart
type Cart struct {
	ID        string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    string     `gorm:"not null" json:"user_id"`
	Items     []CartItem `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE;" json:"items"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// 🧾 CartItem — each product in the cart
type CartItem struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CartID    string    `gorm:"not null" json:"cart_id"`
	ProductID string    `gorm:"not null" json:"product_id"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	Price     float64   `gorm:"default:0" json:"price"` // ✅ New field: store product price at checkout time
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
