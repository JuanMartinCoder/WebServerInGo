package database

import (
	"errors"
	"time"
)

type RefreshToken struct {
	UserId       int       `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	ExpiredAt    time.Time `json:"expired_at"`
}

func (db *DB) SaveRefreshToken(userId int, refreshToken string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	refresh := RefreshToken{
		UserId:       userId,
		RefreshToken: refreshToken,
		ExpiredAt:    time.Now().Add(time.Hour),
	}
	dbStructure.RefreshTokens[refreshToken] = refresh

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetUserByRefreshToken(refreshToken string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	userId, ok := dbStructure.RefreshTokens[refreshToken]
	if !ok {
		return User{}, errors.New("refresh token not found")
	}
	user, ok := dbStructure.Users[userId.UserId]
	if !ok {
		return User{}, errors.New("user not found")
	}
	return user, nil
}

func (db *DB) RevokeRefreshToken(refreshToken string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}
	delete(dbStructure.RefreshTokens, refreshToken)
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}
	return nil
}
