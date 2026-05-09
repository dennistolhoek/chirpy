package main

import (
	"database/sql"
	"errors"
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", nil)
		return
	}
	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if errors.Is(err, sql.ErrNoRows) {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirp", err)
		return
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		UserID:    dbChirp.UserID,
		Body:      dbChirp.Body,
	}

	respondWithJSON(w, http.StatusOK, chirp)

}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	authorID := r.URL.Query().Get("author_id")
	selectSort := r.URL.Query().Get("sort")
	if selectSort != "" && selectSort != "asc" && selectSort != "desc" {
		respondWithError(w, http.StatusBadRequest, "Invalid sort parameter", nil)
		return
	}

	if authorID != "" {
		authorUUID, err := uuid.Parse(authorID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID", nil)
			return
		}
		dbChirps, err := cfg.db.GetChirpsByAuthor(r.Context(), authorUUID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
			return
		}
		chirps := []Chirp{}
		for _, dbChirp := range dbChirps {
			chirps = append(chirps, Chirp{
				ID:        dbChirp.ID,
				CreatedAt: dbChirp.CreatedAt,
				UpdatedAt: dbChirp.UpdatedAt,
				UserID:    dbChirp.UserID,
				Body:      dbChirp.Body,
			})
		}
		sort.Slice(chirps, func(i, j int) bool {
			if selectSort == "asc" {
				return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
			}
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
		respondWithJSON(w, http.StatusOK, chirps)
	} else {
		dbChirps, err := cfg.db.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
			return
		}
		chirps := []Chirp{}
		for _, dbChirp := range dbChirps {
			chirps = append(chirps, Chirp{
				ID:        dbChirp.ID,
				CreatedAt: dbChirp.CreatedAt,
				UpdatedAt: dbChirp.UpdatedAt,
				UserID:    dbChirp.UserID,
				Body:      dbChirp.Body,
			})
		}
		sort.Slice(chirps, func(i, j int) bool {
			if selectSort == "asc" {
				return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
			}
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
		respondWithJSON(w, http.StatusOK, chirps)
	}
}
