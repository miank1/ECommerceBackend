package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads .env file if present
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}
}

// GetEnv returns the environment variable or fallback value
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
