package config

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/michaeldebetaz/chirpy/internal/auth"
	"github.com/michaeldebetaz/chirpy/internal/database"
	"github.com/michaeldebetaz/chirpy/internal/validators"
)

func (c *Config) ChirpGET(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		err := fmt.Errorf("failed to parse chirp ID: %w", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirp, err := c.Queries.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err := fmt.Errorf("chirp not found: %w", err)
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		err := fmt.Errorf("failed to get chirp: %w", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)
}

func (c *Config) ChirpsGET(w http.ResponseWriter, r *http.Request) {
	chirps, err := c.Queries.GetChirps(r.Context())
	if err != nil {
		err := fmt.Errorf("failed to get chirps: %w", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (c *Config) ChirpsPOST(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err := fmt.Errorf("failed to get bearer access token: %w", err)
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userID, err := auth.ValidateJWT(accessToken, c.Env.JWT_SECRET)
	if err != nil {
		err := fmt.Errorf("failed to validate JWT: %w", err)
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	requestBody := validators.ChirpsPOSTRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		err := fmt.Errorf("failed to decode request body: %w", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	data, err := validators.ChirpsPOST(requestBody)
	if err != nil {
		err := fmt.Errorf("failed to validate chirp: %w", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirp, err := c.Queries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   data.Body,
		UserID: userID,
	})
	if err != nil {
		err := fmt.Errorf("failed to create chirp: %w", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}

func (c *Config) LoginPOST(w http.ResponseWriter, r *http.Request) {
	requestBody := validators.LoginPOSTRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		err := fmt.Errorf("failed to decode request body: %w", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resultData, err := validators.LoginPOST(requestBody)

	user, err := c.Queries.GetUserByEmail(r.Context(), resultData.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err := fmt.Errorf("user not found: %w", err)
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		err := fmt.Errorf("failed to get user: %w", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := auth.CheckPasswordHash(user.PasswordHash, resultData.Password); err != nil {
		err := fmt.Errorf("invalid password: %w", err)
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	accessToken, err := auth.GenerateJWT(user.ID, c.Env.JWT_SECRET, resultData.ExpiresInSeconds)
	if err != nil {
		err := fmt.Errorf("failed to generate JWT: %w", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	refreshTokenStr, err := auth.GenerateRefreshToken()
	if err != nil {
		err := fmt.Errorf("failed to generate refresh token: %w", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	refereshToken, err := c.Queries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshTokenStr,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Duration(60 * 24 * time.Hour)),
	})

	responseBody := map[string]any{
		"id":            user.ID,
		"email":         user.Email,
		"token":         accessToken,
		"refresh_token": refereshToken.Token,
		"created_at":    user.CreatedAt,
		"updated_at":    user.UpdatedAt,
	}

	respondWithJSON(w, http.StatusOK, responseBody)
}

func (c *Config) RefreshPOST(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err := fmt.Errorf("failed to get bearer refresh token: %w", err)
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	row, err := c.Queries.GetUserByRefreshToken(r.Context(), refreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err := fmt.Errorf("user not found: %w", err)
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		err := fmt.Errorf("failed to get user: %w", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if row.RefreshToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "refresh token has expired")
		return
	}

	if row.RefreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "refresh token has been revoked")
		return
	}

	expiresIn := time.Duration(60 * 60 * time.Second)
	accessToken, err := auth.GenerateJWT(row.User.ID, c.Env.JWT_SECRET, expiresIn)
	if err != nil {
		err := fmt.Errorf("failed to generate JWT: %w", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseBody := map[string]any{"token": accessToken}
	respondWithJSON(w, http.StatusOK, responseBody)
}

func (c *Config) ResetPOST(w http.ResponseWriter, r *http.Request) {
	if c.Env.PLATFORM != "dev" {
		respondWithError(w, http.StatusForbidden, "Reset is only allowed in development mode")
		return
	}

	if err := c.Queries.DeleteAllUsers(r.Context()); err != nil {
		err := fmt.Errorf("failed to delete all users: %w", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	body := fmt.Sprintf("Hits reset to 0; Users reset to 0\n")
	if _, err := w.Write([]byte(body)); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func (c *Config) RevokeAction(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err := fmt.Errorf("failed to get bearer refresh token: %w", err)
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if err := c.Queries.RevokeRefreshToken(r.Context(), refreshToken); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err := fmt.Errorf("refresh token not found: %w", err)
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		err := fmt.Errorf("failed to revoke refresh token: %w", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusNoContent)
}

func (c *Config) UsersPOST(w http.ResponseWriter, r *http.Request) {
	requestBody := validators.UsersPOSTRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		err := fmt.Errorf("failed to decode request body: %w", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resultData, err := validators.UsersPOST(requestBody)
	if err != nil {
		err := fmt.Errorf("failed to validate email: %w", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := c.Queries.CreateUser(r.Context(), database.CreateUserParams{
		Email:        resultData.Email,
		PasswordHash: resultData.PasswordHash,
	})
	if err != nil {
		err := fmt.Errorf("failed to create user: %w", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseBody := map[string]any{
		"id":         user.ID,
		"email":      user.Email,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}

	respondWithJSON(w, http.StatusCreated, responseBody)
}

func (c *Config) UsersPUT(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err := fmt.Errorf("failed to get bearer access token: %w", err)
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	requestBody := validators.UsersPUTRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		err := fmt.Errorf("failed to decode request body: %w", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resultData, err := validators.UsersPUT(requestBody)
	if err != nil {
		err := fmt.Errorf("failed to validate email: %w", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID, err := auth.ValidateJWT(accessToken, c.Env.JWT_SECRET)
	if err != nil {
		err := fmt.Errorf("failed to validate JWT: %w", err)
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	user, err := c.Queries.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:        resultData.Email,
		PasswordHash: resultData.PasswordHash,
		ID:           userID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err := fmt.Errorf("user not found: %w", err)
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		err := fmt.Errorf("failed to update user: %w", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseBody := map[string]any{
		"id":         user.ID,
		"email":      user.Email,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}
	respondWithJSON(w, http.StatusOK, responseBody)
}

func respondWithJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		err := fmt.Errorf("failed to encode response: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, map[string]string{"error": message})
}
