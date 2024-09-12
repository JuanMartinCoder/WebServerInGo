package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) handlePostRefresh(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Token string `json:"token"`
	}
	headerAuth := r.Header.Get("Authorization")
	authToken, ok := strings.CutPrefix(headerAuth, "Bearer ")
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "No token provided")
		return
	}

	user, err := cfg.DB.GetUserByRefreshToken(authToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token not found")
		return
	}

	claim := jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", user.ID),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		Issuer:    "chirpy",
	}
	tokenC := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	token, err := tokenC.SignedString([]byte(cfg.jwtSecret))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to sign token")
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		Token: token,
	})
}
