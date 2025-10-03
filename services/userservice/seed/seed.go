package main

import (
	"ecommerce-backend/pkg/db"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"uniqueIndex"`
	Password string `json:"password"`
}

func main() {
	cfg := db.Config{
		DSN:         os.Getenv("DATABASE_DSN"),
		MaxRetries:  6,
		RetryDelay:  2 * time.Second,
		ConnTimeout: 5 * time.Second,
	}

	gormDB, err := db.InitPostgres(cfg)
	if err != nil {
		log.Fatalf("could not initialize database: %v", err)
	}

	if err := gormDB.AutoMigrate(&User{}); err != nil {
		log.Fatalf("✅ Migration failed: %v", err)
	}

	// Clear old data
	if err := gormDB.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE").Error; err != nil {
		log.Fatalf("Failed to truncate users table: %v", err)
	}
	log.Println("🧹 Cleared existing users")

	// Load JSON
	file, err := os.Open("/seed/users.json")
	if err != nil {
		log.Fatalf("Failed to open JSON: %v", err)
	}
	defer file.Close()

	var users []User
	if err := json.NewDecoder(file).Decode(&users); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	// Hash passwords
	for i := range users {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(users[i].Password), bcrypt.DefaultCost)
		users[i].Password = string(hashed)
	}

	// Insert into DB
	if err := gormDB.Create(&users).Error; err != nil {
		log.Fatalf("Failed to insert users: %v", err)
	}

	fmt.Printf("✅ Inserted %d users into UserDB\n", len(users))
}
