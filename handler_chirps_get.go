package main

import (
	"learn_http_servers/internal/database"
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func authorIDFromRequest(r *http.Request) (uuid.UUID, error) {
	author := r.URL.Query().Get("author_id")
	if author == "" {
		return uuid.Nil, nil
	}
	authorUUID, err := uuid.Parse(author)
	if err != nil {
		return uuid.Nil, err
	}
	return authorUUID, nil
}

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	authorID, err := authorIDFromRequest(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
		return
	}

	var dbChirps []database.Chirp

	if authorID != uuid.Nil {
		dbChirps, err = cfg.db.GetChirpsByUserID(r.Context(), authorID)
	} else {
		dbChirps, err = cfg.db.GetAllChirps(r.Context())
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
		return
	}

	chirps := []Chirp{}
	for _, chirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	if r.URL.Query().Get("sort") == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
