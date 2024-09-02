package main

import (
	"fmt"
	"net/http"
)

func main() {
	serMux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: serMux,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Print(err)
	}
}
