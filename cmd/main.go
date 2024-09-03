package main

import (
	"fmt"
	"net/http"
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
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", c.fileserverHits)))
}

func (c *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	c.fileserverHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func main() {
	apiCfg := apiConfig{
		fileserverHits: 0,
	}
	serMux := http.NewServeMux()
	serMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir("../ui/")))))
	serMux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	serMux.HandleFunc("GET /metrics", apiCfg.handlerMetrics)
	serMux.HandleFunc("/reset", apiCfg.handleReset)

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
