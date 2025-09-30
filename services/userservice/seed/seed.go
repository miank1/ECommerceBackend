package main

import (
	"ecommerce-backend/pkg/config"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"uniqueIndex"`
	Password string `json:"password"`
}

func main() {
	config.LoadEnv()
	dsn := os.Getenv("DATABASE_DSN") // e.g., "host=postgres user=root password=secret dbname=ecommerce port=5432 sslmode=disable"
	log.Fatal("Hello World DATABASE_DSN not set !!", dsn)

	if dsn == "" {
		log.Fatal("DATABASE_DSN not set!")
	} else {
		log.Println("✅ DATABASE_DSN loaded:", dsn)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("✅ Failed to connect DB: %v", err)
	}

	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatalf("✅ Migration failed: %v", err)
	}

	// Clear old data
	if err := db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE").Error; err != nil {
		log.Fatalf("Failed to truncate users table: %v", err)
	}
	log.Println("🧹 Cleared existing users")

	// Load JSON
	file, err := os.Open("seed/users.json")
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
	if err := db.Create(&users).Error; err != nil {
		log.Fatalf("Failed to insert users: %v", err)
	}

	fmt.Printf("✅ Inserted %d users into UserDB\n", len(users))
}
