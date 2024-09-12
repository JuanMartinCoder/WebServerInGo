package main

import (
	"net/http"

	"github.com/JuanMartinCoder/WebServerInGo/internal/auth"
)

func (cfg *apiConfig) handlePostRevoke(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Token string `json:"token"`
	}

	authToken, err := auth.GetHeaderBearer(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No token provided")
		return
	}

	err = cfg.DB.RevokeRefreshToken(authToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	respondWithJSON(w, http.StatusNoContent, Response{Token: authToken})
}
