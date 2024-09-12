package database

import (
	"errors"
	"strconv"
)

type Chirp struct {
	ID       int    `json:"id"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
}

func (db *DB) GetChirpById(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	chirp, ok := dbStructure.Chrips[id]
	if !ok {
		return Chirp{}, errors.New("chirp not found")
	}
	return chirp, nil
}

func (db *DB) CreateChirp(body string, userId string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chrips) + 1
	chirp := Chirp{
		ID:       id,
		Body:     body,
		AuthorId: userIdInt,
	}

	dbStructure.Chrips[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}
	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	chirps := make([]Chirp, 0, len(dbStructure.Chrips))
	for _, chirp := range dbStructure.Chrips {
		chirps = append(chirps, chirp)
	}
	return chirps, nil
}

func (db *DB) DeleteChirp(id int, userId string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return err
	}

	if userIdInt != dbStructure.Chrips[id].AuthorId {
		return errors.New("you are not the author of this chirp")
	}

	delete(dbStructure.Chrips, id)
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}
	return nil
}
