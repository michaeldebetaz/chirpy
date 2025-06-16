package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	password := "correcthorsebatterystaple"
	hash, err := Hash(password)
	if err != nil {
		err = fmt.Errorf("failed to hash password '%s': %w", password, err)
		t.Fatal(err.Error())
	}

	if hash == "" {
		t.Fatal("expected a non-empty hashed password")
	}

	if err = CheckPasswordHash(hash, password); err != nil {
		err := fmt.Errorf("expected passwords to match: %w", err)
		t.Fatal(err.Error())
	}
}

func TestJWTGenerationAndValidation(t *testing.T) {
	tokenSecret := "4CntggtrGc3YEsKEUA2MEXcc/HiCp8J5Y5oec/jrpNhEtkw6HvmOSXxUA7rY3Id1ZgEq/NTVzFVtEZwYhyAw3g="

	userID := uuid.New()
	token, err := GenerateJWT(userID, tokenSecret, time.Duration(24*time.Hour))
	if err != nil {
		t.Fatalf("failed to generate JWT: %v", err)
	}

	parsedUserID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("failed to validate JWT: %v", err)
	}

	if parsedUserID != userID {
		t.Fatalf("expected user ID %s, got %s", userID, parsedUserID)
	}
}
