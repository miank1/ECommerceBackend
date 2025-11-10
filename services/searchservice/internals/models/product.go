package models

import uuid "github.com/google/uuid"

type Product struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name     string    `json:"name"`
	Category string    `json:"category"`
	Price    float64   `json:"price"`
	Stock    int       `json:"stock"`
}
