package validators

import (
	"time"
)

type ChirpsActionRequestBody struct {
	Body string `json:"body"`
}

type ChirpsActionResultData struct {
	Body string `json:"body"`
}

type LoginActionRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginActionResultData struct {
	Email            string        `json:"email"`
	Password         string        `json:"password"`
	ExpiresInSeconds time.Duration `json:"expires_in_seconds"`
}

type UsersActionRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UsersActionResultData struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}
