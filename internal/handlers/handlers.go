package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/michaeldebetaz/chirpy/internal/database"
	"github.com/michaeldebetaz/chirpy/internal/middlewares"
	"github.com/michaeldebetaz/chirpy/internal/state"
	"github.com/michaeldebetaz/chirpy/internal/validators"
)

func Healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	body := fmt.Sprintf("%s", http.StatusText(http.StatusOK))
	w.Write([]byte(body))
}

func Metrics(w http.ResponseWriter, r *http.Request) {
	hits, ok := r.Context().Value(middlewares.HITS_KEY).(int32)
	if !ok {
		http.Error(w, "Failed to retrieve hits from context", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	body := fmt.Sprintf(`
<html> 
	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p> 
	</body>
</html>`, hits)
	if _, err := w.Write([]byte(body)); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func Reset(s *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.Env.PLATFORM != "dev" {
			respondWithError(w, http.StatusForbidden, "Reset is only allowed in development mode")
			return
		}

		if err := s.Queries.DeleteAllUsers(r.Context()); err != nil {
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
}

func Users(s *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type RequestBody struct {
			Email string `json:"email"`
		}

		requestBody := RequestBody{}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			err := fmt.Errorf("failed to decode request body: %w", err)
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		email, err := validators.Email(requestBody.Email)
		if err != nil {
			err := fmt.Errorf("failed to validate email: %w", err)
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		user, err := s.Queries.CreateUser(r.Context(), email)
		if err != nil {
			err := fmt.Errorf("failed to create user: %w", err)
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusCreated, user)
	}
}

func ChirpLoader(s *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chirpID, err := validators.UUID(r.PathValue("chirpID"))
		if err != nil {
			err := fmt.Errorf("failed to parse chirp ID: %w", err)
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		chirp, err := s.Queries.GetChirpByID(r.Context(), chirpID)
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
}

func ChirpsLoader(s *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chirps, err := s.Queries.GetChirps(r.Context())
		if err != nil {
			err := fmt.Errorf("failed to get chirps: %w", err)
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, chirps)
	}
}

func ChirpsAction(s *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestBody := validators.ChirpsActionRequestBody{}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			err := fmt.Errorf("failed to decode request body: %w", err)
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		data, err := validators.ChirpsAction(requestBody)
		if err != nil {
			err := fmt.Errorf("failed to validate chirp: %w", err)
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		user, err := s.Queries.CreateChirp(r.Context(), database.CreateChirpParams{
			Body:   data.Body,
			UserID: data.UserID,
		})
		if err != nil {
			err := fmt.Errorf("failed to create chirp: %w", err)
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusCreated, user)
	}
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
