package models

import "time"

type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
	Price       float64   `json:"price" validate:"required,gt=0"`
	Stock       int       `json:"stock" validate:"required,gte=0"`
	CategoryID  string    `json:"category_id" validate:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductResponse struct {
	Status  string   `json:"status"`
	Message string   `json:"message,omitempty"`
	Data    *Product `json:"data,omitempty"`
}

type ProductListResponse struct {
	Status  string     `json:"status"`
	Message string     `json:"message,omitempty"`
	Data    []Product  `json:"data,omitempty"`
}
