package main

import (
	"net/http"

	"github.com/JuanMartinCoder/WebServerInGo/internal/auth"
)

func (cfg *apiConfig) handlePostRefresh(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Token string `json:"token"`
	}
	authToken, err := auth.GetHeaderBearer(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No token provided")
		return
	}

	user, err := cfg.DB.GetUserByRefreshToken(authToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token not found")
		return
	}
	token, err := auth.CreateAccessToken(user.ID, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to sign token")
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		Token: token,
	})
}
