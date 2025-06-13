package validators

import "github.com/google/uuid"

type ChirpsActionRequestBody struct {
	Body   string `json:"body"`
	UserID string `json:"user_id"`
}

type ChirpsActionResultData struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}
