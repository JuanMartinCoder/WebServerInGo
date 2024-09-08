package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type apiConfig struct {
	fileserverHits int
}

func (c *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (c *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("<html> <body>  <h1>Welcome, Chirpy Admin</h1>   <p>Chirpy has been visited %d times!</p> </body></html>", c.fileserverHits)))
}

func (c *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	c.fileserverHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}
	type successResp struct {
		CleanBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	chirData := chirp{}
	err := decoder.Decode(&chirData)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	if len(chirData.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	respondWithJSON(w, http.StatusOK, successResp{
		CleanBody: cleanBody(chirData.Body),
	})
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

func main() {
	apiCfg := apiConfig{
		fileserverHits: 0,
	}
	serMux := http.NewServeMux()
	serMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir("../ui/")))))
	serMux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	serMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	serMux.HandleFunc("/api/reset", apiCfg.handleReset)
	serMux.HandleFunc("POST /api/validate_chirp", handleValidateChirp)

	server := &http.Server{
		Addr:    ":8080",
		Handler: serMux,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println("Server started at port: ", server.Addr)
}
