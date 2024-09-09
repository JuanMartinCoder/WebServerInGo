package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func (cfg *apiConfig) handleGetChirpID(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("chirpId")
	valueId, err := strconv.Atoi(path)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse chirp id")
		return
	}
	chirp, err := cfg.DB.GetChirpById(valueId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)
}

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	chirpsDb, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps")
		return
	}
	chirps := []Chirp{}
	for _, dbchirp := range chirpsDb {
		chirps = append(chirps, Chirp{
			ID:   dbchirp.ID,
			Body: dbchirp.Body,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlePostChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	chirData := parameters{}
	err := decoder.Decode(&chirData)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	cleanedBody, err := validateChirp(chirData.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirp, err := cfg.DB.CreateChirp(cleanedBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:   chirp.ID,
		Body: chirp.Body,
	})
}

func validateChirp(body string) (string, error) {
	if len(body) > 140 {
		return "", errors.New("Chirp is too long")
	}

	cleanBody := cleanBody(body)
	return cleanBody, nil
}

func cleanBody(body string) string {
	type badwords struct {
		words []string
	}
	words := badwords{
		words: []string{
			"kerfuffle",
			"sharbert",
			"fornax",
		},
	}

	splitedBody := strings.Split(body, " ")
	for i, word := range splitedBody {
		for _, badword := range words.words {
			if strings.ToLower(word) == badword {
				splitedBody[i] = "****"
				break
			}
		}
	}
	return strings.Join(splitedBody, " ")
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}
