package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Product struct {
	ID       uint    `gorm:"primaryKey"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
	Stock    int     `json:"stock"`
}

func main() {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("DATABASE_DSN not set in environment")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}

	if err := db.AutoMigrate(&Product{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Reset products
	if err := db.Exec("TRUNCATE TABLE products RESTART IDENTITY CASCADE").Error; err != nil {
		log.Fatalf("Failed to truncate table: %v", err)
	}
	log.Println("🧹 Cleared existing products")

	// Open JSON file
	file, err := os.Open("seed/products.json")
	if err != nil {
		log.Fatalf("Failed to open JSON file: %v", err)
	}
	defer file.Close()

	var products []Product
	if err := json.NewDecoder(file).Decode(&products); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	// Insert products
	if err := db.Create(&products).Error; err != nil {
		log.Fatalf("Failed to insert products: %v", err)
	}

	fmt.Printf("✅ Inserted %d products from JSON\n", len(products))
}
