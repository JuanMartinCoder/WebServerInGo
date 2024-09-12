package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.CreateLogin(params.Email, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Access Token
	claim := jwt.RegisteredClaims{}
	claim = jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(1 * time.Hour)),
		Subject:   fmt.Sprintf("%d", user.ID),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(cfg.jwtSecret))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't sign token")
		return
	}
	// Refresh Token
	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate random bytes")
		return
	}
	refreshToken := hex.EncodeToString(b)

	err = cfg.DB.SaveRefreshToken(user.ID, refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:           user.ID,
		Email:        user.Email,
		Token:        tokenString,
		RefreshToken: refreshToken,
	})
}
