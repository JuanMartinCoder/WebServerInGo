package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/JuanMartinCoder/WebServerInGo/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	jwtSecret      string
	PolkaToken     string
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

func main() {
	// Load environment variables
	godotenv.Load("../.env")

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
		jwtSecret:      os.Getenv("JWT_SECRET"),
		PolkaToken:     os.Getenv("POLKA_TOKEN"),
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
	serMux.HandleFunc("POST /api/chirps", apiCfg.handlePostChirp)
	serMux.HandleFunc("GET /api/chirps", apiCfg.handleGetChirps)
	serMux.HandleFunc("GET /api/chirps/{chirpId}", apiCfg.handleGetChirpID)
	serMux.HandleFunc("DELETE /api/chirps/{chirpId}", apiCfg.handleDeleteChirp)

	serMux.HandleFunc("POST /api/users", apiCfg.handlePostUsers)
	serMux.HandleFunc("PUT /api/users", apiCfg.handlePutUsers)

	serMux.HandleFunc("POST /api/login", apiCfg.handleLogin)
	serMux.HandleFunc("POST /api/refresh", apiCfg.handlePostRefresh)
	serMux.HandleFunc("POST /api/revoke", apiCfg.handlePostRevoke)

	serMux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlePostPolkaWebhooks)

	server := &http.Server{
		Addr:    ":8080",
		Handler: serMux,
	}

	error := server.ListenAndServe()
	if err != nil {
		fmt.Print(error)
	}
	fmt.Println("Server started at port: ", server.Addr)
}
