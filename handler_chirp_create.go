package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	auth "learn_http_servers/internal/auth"
	"learn_http_servers/internal/database"

	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type response struct {
		Chirp
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithJSON(w, http.StatusUnauthorized, "Invalid authorization header")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.serverSecret)
	if err != nil {
		respondWithJSON(w, http.StatusUnauthorized, "Invalid JWT token")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding JSON parameters", err)
		return
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleaned,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		Chirp: Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      params.Body,
			UserID:    userID,
		},
	})
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140

	if len(body) > maxChirpLength {
		return "", errors.New("chirp is too long")
	}

	cleaned := badwordCheck(body)
	return cleaned, nil
}
