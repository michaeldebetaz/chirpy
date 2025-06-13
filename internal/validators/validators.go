package validators

import (
	"fmt"
	"net/mail"
	"strings"

	"github.com/google/uuid"
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

	userID, err := uuid.Parse(chirp.UserID)
	if err != nil {
		err := fmt.Errorf("failed to parse user ID '%s': %w", chirp.UserID, err)
		return nil, err
	}

	return &ChirpsActionResultData{Body: body, UserID: userID}, nil
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

func Email(e string) (string, error) {
	address, err := mail.ParseAddress(e)
	if err != nil {
		err := fmt.Errorf("failed to parse email address '%s': %w", e, err)
		return "", err
	}

	email := address.Address

	return email, nil
}
