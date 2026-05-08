package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dennistolhoek/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		UserID string `json:"user_id"`
	}
	type parameters struct {
		Event string `json:"event"`
		Data  Data   `json:"data"`
	}
	var params parameters

	apiKey, err := auth.GetAPIkey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid API key", err)
		return
	}
	if apiKey != cfg.polkaSecret {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key", nil)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Invalid user ID", err)
		return
	}

	err = cfg.db.UpdateUserChirpyRed(r.Context(), userID)
	if errors.Is(err, sql.ErrNoRows) {
		respondWithError(w, http.StatusNotFound, "User not found", err)
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
