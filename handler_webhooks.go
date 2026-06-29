package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpdateRedStatus(w http.ResponseWriter, r *http.Request) {
	type Parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		}
	}

	decoder := json.NewDecoder(r.Body)
	params := Parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding JSON parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.db.UpdateRedStatus(r.Context(), params.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "User not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user status", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
