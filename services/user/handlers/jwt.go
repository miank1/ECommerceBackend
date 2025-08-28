package handlers

import (
	"ecommerce-api/services/user/models"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	// In a real application, this should be in a secure configuration
	jwtSecret = "your-256-bit-secret"
	// Token expires in 24 hours
	tokenExpiration = 24 * time.Hour
)

// Claims represents the JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

// generateToken creates a new JWT token for a user
func generateToken(user models.User) (string, error) {
	// Create the expiration time
	expirationTime := time.Now().Add(tokenExpiration)

	// Create the JWT claims
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   user.ID,
		},
	}

	// Generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string
	return token.SignedString([]byte(jwtSecret))
}
