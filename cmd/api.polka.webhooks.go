package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlePostPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	type PolkaEvent struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}

	headerAuth := r.Header.Get("Authorization")
	authToken, ok := strings.CutPrefix(headerAuth, "ApiKey ")
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Invalid authorization header")
		return
	}

	if authToken != cfg.PolkaToken {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key")
		return
	}

	params := PolkaEvent{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "Invalid event")
		return
	}

	_, err = cfg.DB.UpdateUserChirpyRed(params.Data.UserID, true)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found")
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
