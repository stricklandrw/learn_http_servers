package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpGet(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to parse chirp UUID", err)
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if dbChirp.ID == uuid.Nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", nil)
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirp", err)
		return
	}

	type response struct {
		Chirp
	}

	respondWithJSON(w, http.StatusOK, response{
		Chirp: Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		},
	})
}
