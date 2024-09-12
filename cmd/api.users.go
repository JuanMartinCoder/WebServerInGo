package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (cfg *apiConfig) handlePostUsers(w http.ResponseWriter, r *http.Request) {
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

	user, err := cfg.DB.CreateUser(params.Email, params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:    user.ID,
		Email: user.Email,
	})
}

func (cfg *apiConfig) handlePutUsers(w http.ResponseWriter, r *http.Request) {
	// Get Params and bearer token
	type Parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	headerAuth := r.Header.Get("Authorization")
	authToken, ok := strings.CutPrefix(headerAuth, "Bearer ")
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "No token provided")
		return
	}
	claim := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(authToken, &claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.jwtSecret), nil
	})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}
	claims := token.Claims
	userId, err := claims.GetSubject()
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid Claims")
		return
	}
	decoder := json.NewDecoder(r.Body)
	params := Parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	updatedUser, err := cfg.DB.UpdateUser(userId, params.Email, params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:    updatedUser.ID,
		Email: updatedUser.Email,
	})
}
