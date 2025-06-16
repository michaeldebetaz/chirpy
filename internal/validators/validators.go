package validators

import (
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/michaeldebetaz/chirpy/internal/auth"
)

func BadWords() map[string]bool {
	return map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}
}

func ChirpsPOST(chirp ChirpsPOSTRequestBody) (*ChirpsPOSTResultData, error) {
	body := strings.TrimSpace(chirp.Body)

	if len(body) > 140 {
		err := fmt.Errorf("chrip length (%d) exceeds 140 characters", len(body))
		return nil, err
	}

	words := strings.Fields(body)
	badWords := BadWords()

	for i, word := range words {
		w := strings.ToLower(word)
		if _, ok := badWords[w]; ok {
			words[i] = "****"
		}
	}

	body = strings.Join(words, " ")

	return &ChirpsPOSTResultData{Body: body}, nil
}

func LoginPOST(requestBody LoginPOSTRequestBody) (*LoginPOSTResultData, error) {
	emailAddress, err := mail.ParseAddress(requestBody.Email)
	if err != nil {
		err := fmt.Errorf("failed to parse email address '%s': %w", requestBody.Email, err)
		return nil, err
	}

	if requestBody.Password == "" {
		err := fmt.Errorf("password cannot be empty")
		return nil, err
	}

	return &LoginPOSTResultData{
		Email:            emailAddress.Address,
		Password:         requestBody.Password,
		ExpiresInSeconds: time.Duration(60 * 60 * time.Second),
	}, nil
}

func UsersPOST(requestBody UsersPOSTRequestBody) (*UsersPOSTResultData, error) {
	emailAddress, err := mail.ParseAddress(requestBody.Email)
	if err != nil {
		err := fmt.Errorf("failed to parse email address '%s': %w", requestBody.Email, err)
		return nil, err
	}

	if requestBody.Password == "" {
		err := fmt.Errorf("password cannot be empty")
		return nil, err
	}

	passwordHash, err := auth.Hash(requestBody.Password)
	if err != nil {
		err := fmt.Errorf("failed to hash password: %w", err)
		return nil, err
	}

	return &UsersPOSTResultData{
		Email:        emailAddress.Address,
		PasswordHash: passwordHash,
	}, nil
}

func UsersPUT(requestBody UsersPUTRequestBody) (*UsersPUTResultData, error) {
	if requestBody.Email == "" {
		err := fmt.Errorf("email cannot be empty")
		return nil, err
	}

	emailAddress, err := mail.ParseAddress(requestBody.Email)
	if err != nil {
		err := fmt.Errorf("failed to parse email address '%s': %w", requestBody.Email, err)
		return nil, err
	}

	if requestBody.Password == "" {
		err := fmt.Errorf("password cannot be empty")
		return nil, err
	}

	passwordHash, err := auth.Hash(requestBody.Password)
	if err != nil {
		err := fmt.Errorf("failed to hash password: %w", err)
		return nil, err
	}

	return &UsersPUTResultData{
		Email:        emailAddress.Address,
		PasswordHash: passwordHash,
	}, nil
}
