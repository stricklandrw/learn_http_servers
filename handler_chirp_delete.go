package main

import (
	"net/http"

	"learn_http_servers/internal/auth"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	if chirpIDString == "" {
		respondWithJSON(w, http.StatusBadRequest, "Chirp ID is required")
		return
	}

	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, "Invalid chirp ID")
		return
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

	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirp", err)
		return
	}

	if chirp.UserID != userID {
		respondWithJSON(w, http.StatusForbidden, "You are not authorized to delete this chirp")
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
