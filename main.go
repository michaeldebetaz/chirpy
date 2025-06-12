package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/michaeldebetaz/chirpy/internal/handlers"
	"github.com/michaeldebetaz/chirpy/internal/middlewares"
)

func main() {
	mux := http.ServeMux{}

	cfg := middlewares.NewConfig()

	fileServerHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	mux.Handle("GET /app/", cfg.IncrementFileserverHits(fileServerHandler))
	mux.HandleFunc("GET /api/healthz", handlers.Healthz)
	mux.Handle("GET /api/metrics", cfg.FileserverHits())
	mux.Handle("POST /api/reset", cfg.ResetFileserverHits())

	server := &http.Server{
		Handler: &mux,
		Addr:    ":8080",
	}

	fmt.Printf("Starting server on %s\n", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
