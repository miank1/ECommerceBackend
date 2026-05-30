package jwtutil

import (
	"os"
	"testing"
)

func TestGenerateToken(t *testing.T) {

	os.Setenv("JWT_SECRET", "test-secret")

	token, err := GenerateToken("user-123")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token == "" {
		t.Fatal("expected token to be generated")
	}
}

func TestParseToken_Success(t *testing.T) {

	secret := "test-secret"
	os.Setenv("JWT_SECRET", secret)

	token, err := GenerateToken("user-123")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	claims, err := ParseToken(secret, token)
	if err != nil {
		t.Fatalf("expected valid token, got error: %v", err)
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		t.Fatal("user_id claim missing")
	}

	if userID != "user-123" {
		t.Fatalf(
			"expected user_id=user-123, got %s",
			userID,
		)
	}
}

func TestParseToken_InvalidToken(t *testing.T) {

	secret := "test-secret"

	_, err := ParseToken(secret, "invalid-token")

	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}
