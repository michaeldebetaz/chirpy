package validators

import (
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/michaeldebetaz/chirpy/internal/auth"
)

func BadWords() map[string]bool {
	return map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}
}

func ChirpsAction(chirp ChirpsActionRequestBody) (*ChirpsActionResultData, error) {
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

	return &ChirpsActionResultData{Body: body}, nil
}

func LoginAction(requestBody LoginActionRequestBody) (*LoginActionResultData, error) {
	emailAddress, err := mail.ParseAddress(requestBody.Email)
	if err != nil {
		err := fmt.Errorf("failed to parse email address '%s': %w", requestBody.Email, err)
		return nil, err
	}

	if requestBody.Password == "" {
		err := fmt.Errorf("password cannot be empty")
		return nil, err
	}

	return &LoginActionResultData{
		Email:            emailAddress.Address,
		Password:         requestBody.Password,
		ExpiresInSeconds: time.Duration(60 * 60 * time.Second),
	}, nil
}

func UsersAction(requestBody UsersActionRequestBody) (*UsersActionResultData, error) {
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

	return &UsersActionResultData{
		Email:        emailAddress.Address,
		PasswordHash: passwordHash,
	}, nil
}

func UUID(u string) (uuid.UUID, error) {
	if u == "" {
		err := fmt.Errorf("UUID cannot be empty")
		return uuid.Nil, err
	}

	id, err := uuid.Parse(u)
	if err != nil {
		err := fmt.Errorf("failed to parse UUID '%s': %w", u, err)
		return uuid.Nil, err
	}

	return id, nil
}
