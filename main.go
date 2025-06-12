package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/michaeldebetaz/chirpy/internal/handlers"
)

func main() {
	handler := http.ServeMux{}

	fileServer := http.FileServer(http.Dir("."))
	handler.Handle("GET /app/", http.StripPrefix("/app/", fileServer))
	handler.HandleFunc("GET /healthz", handlers.Healthz)

	server := &http.Server{
		Handler: &handler,
		Addr:    ":8080",
	}

	fmt.Printf("Starting server on %s\n", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
