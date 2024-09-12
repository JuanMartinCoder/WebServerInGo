package main

import (
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlePostRevoke(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Token string `json:"token"`
	}
	headerAuth := r.Header.Get("Authorization")
	authToken, ok := strings.CutPrefix(headerAuth, "Bearer ")
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "No token provided")
		return
	}

	err := cfg.DB.RevokeRefreshToken(authToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	respondWithJSON(w, http.StatusNoContent, Response{Token: authToken})
}
