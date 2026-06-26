package main

import (
	"encoding/json"
	"learn_http_servers/internal/auth"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Expiry   int    `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding JSON parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Incorrect email or password", err)
		return
	}

	isMatch, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !isMatch {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	tokenDuration := time.Hour
	if params.Expiry > 0 && params.Expiry < 3600 {
		tokenDuration = time.Duration(params.Expiry) * time.Second
	}

	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.serverSecret,
		tokenDuration,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: accessToken,
	})
}
