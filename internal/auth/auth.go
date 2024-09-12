package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	refreshToken := hex.EncodeToString(b)
	return refreshToken, nil
}

func CreateAccessToken(id int, jwtSecret string) (string, error) {
	claim := jwt.RegisteredClaims{}
	claim = jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(1 * time.Hour)),
		Subject:   fmt.Sprintf("%d", id),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateAccessToken(authToken string, jwtSecret string) (string, error) {
	claim := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(authToken, &claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return "", err
	}
	claims := token.Claims
	userId, err := claims.GetSubject()
	if err != nil {
		return "", err
	}
	return userId, nil
}

func GetHeaderBearer(r *http.Request) (string, error) {
	headerAuth := r.Header.Get("Authorization")
	authToken, ok := strings.CutPrefix(headerAuth, "Bearer ")
	if !ok {
		return "", fmt.Errorf("invalid header")
	}
	return authToken, nil
}
