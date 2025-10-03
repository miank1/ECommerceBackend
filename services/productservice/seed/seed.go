package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	uuid "github.com/jackc/pgx/pgtype/ext/satori-uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Product struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name     string    `json:"name"`
	Category string    `json:"category"`
	Price    float64   `json:"price"`
	Stock    int       `json:"stock"`
}

func main() {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("DATABASE_DSN not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	if err := db.AutoMigrate(&Product{}); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	// clear existing
	db.Exec("TRUNCATE TABLE products RESTART IDENTITY CASCADE")

	// open JSON file
	file, err := os.Open("/seed/products.json")
	if err != nil {
		log.Fatalf("failed to open JSON: %v", err)
	}
	defer file.Close()

	var products []Product
	if err := json.NewDecoder(file).Decode(&products); err != nil {
		log.Fatalf("failed to decode JSON: %v", err)
	}

	// insert
	if err := db.Create(&products).Error; err != nil {
		log.Fatalf("failed to insert: %v", err)
	}

	fmt.Printf("✅ Inserted %d products\n", len(products))
}
