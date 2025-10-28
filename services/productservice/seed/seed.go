package seed

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"os"

	uuid "github.com/jackc/pgx/pgtype/ext/satori-uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name     string    `json:"name"`
	Category string    `json:"category"`
	Price    float64   `json:"price"`
	Stock    int       `json:"stock"`
}

// SeedProducts loads products.json and inserts into DB
func SeedProducts(db *gorm.DB) {
	start := time.Now()

	data, err := os.ReadFile("../seed/products.json")
	if err != nil {
		log.Fatalf("❌ Failed to read products.json: %v", err)
	}

	var products []Product
	if err := json.Unmarshal(data, &products); err != nil {
		log.Fatalf("❌ Failed to parse products.json: %v", err)
	}

	// ✅ Fetch all existing product names in one query
	var existingNames []string
	if err := db.Model(&Product{}).Select("name").Find(&existingNames).Error; err != nil {
		log.Fatalf("❌ Failed to fetch existing product names: %v", err)
	}

	existingSet := make(map[string]bool, len(existingNames))
	for _, name := range existingNames {
		existingSet[name] = true
	}

	// ✅ Collect only new products
	var newProducts []Product
	for _, p := range products {
		if !existingSet[p.Name] {
			newProducts = append(newProducts, p)
		}
	}

	// ✅ Batch insert all new records at once
	if len(newProducts) > 0 {
		if err := db.CreateInBatches(newProducts, 100).Error; err != nil {
			log.Fatalf("❌ Failed to batch insert products: %v", err)
		}
		fmt.Printf("✅ Inserted %d new products (%.2fs)\n", len(newProducts), time.Since(start).Seconds())
	} else {
		fmt.Println("ℹ️ All products already exist — nothing to insert.")
	}
}
