package database

import (
	"errors"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}

func (db *DB) CreateUser(email string, pass string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(pass), 14)

	id := len(dbStructure.Users) + 1
	user := User{
		ID:       id,
		Email:    email,
		Password: hashedPass,
	}
	for _, user := range dbStructure.Users {
		if user.Email == email {
			return User{}, errors.New("email already exists")
		}
	}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (db *DB) UpdateUser(userId string, newEmail string, newPass string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(newPass), 14)

	userIdnew, err := strconv.Atoi(userId)
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[userIdnew]
	if !ok {
		return User{}, errors.New("user not found")
	}
	user.ID = userIdnew
	user.Email = newEmail
	user.Password = hashedPass

	dbStructure.Users[userIdnew] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (db *DB) CreateLogin(email string, pass string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	userData := User{}
	for _, user := range dbStructure.Users {
		if user.Email == email {
			err := bcrypt.CompareHashAndPassword(user.Password, []byte(pass))
			if err != nil {
				return User{}, errors.New("Password or Eamil doesn't match")
			}
			userData.ID = user.ID
			userData.Email = user.Email
			return userData, nil
		}
	}
	return User{}, nil
}
