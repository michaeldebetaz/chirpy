package validators

import (
	"time"
)

type ChirpsPOSTRequestBody struct {
	Body string `json:"body"`
}

type ChirpsPOSTResultData struct {
	Body string `json:"body"`
}

type LoginPOSTRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginPOSTResultData struct {
	Email            string        `json:"email"`
	Password         string        `json:"password"`
	ExpiresInSeconds time.Duration `json:"expires_in_seconds"`
}

type UsersPOSTRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UsersPOSTResultData struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

type UsersPUTRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UsersPUTResultData struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}
