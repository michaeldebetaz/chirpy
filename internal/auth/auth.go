package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CheckPasswordHash(hash string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return fmt.Errorf("password does not match: %w", err)
	}

	return nil
}

func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		err := fmt.Errorf("failed to hash password: %w", err)
		return "", err
	}

	return string(bytes), nil
}

func GetPolkaKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing Authorization header")
	}

	apiKey := strings.TrimPrefix(authHeader, "ApiKey ")

	return apiKey, nil
}

func GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		err := fmt.Errorf("failed to generate random bytes for refresh token: %w", err)
		return "", err
	}

	token := hex.EncodeToString(bytes)

	return token, nil
}

// JWT

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing Authorization header")
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	return token, nil
}

func GenerateJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	s, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		err := fmt.Errorf("failed to create a signed JWT: %w", err)
		return "", err
	}

	return s, nil
}

func ValidateJWT(tokenString string, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse JWT: %w", err)
	}

	if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid JWT")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse user ID from JWT: %w", err)
	}

	return userID, nil
}
