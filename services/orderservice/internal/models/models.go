package models

import "time"

type Order struct {
	ID         string      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID     string      `gorm:"not null" json:"user_id"`
	Status     string      `gorm:"default:'pending'" json:"status"`
	TotalPrice float64     `json:"total_price"`
	Items      []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE;" json:"items"`
	CreatedAt  time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}

type OrderItem struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	OrderID   string    `gorm:"not null" json:"order_id"`
	ProductID string    `gorm:"not null" json:"product_id"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	Price     float64   `gorm:"not null" json:"price"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
