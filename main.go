package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/michaeldebetaz/chirpy/internal/handlers"
	"github.com/michaeldebetaz/chirpy/internal/state"

	_ "github.com/lib/pq"
)

func main() {
	s, err := state.Init()
	if err != nil {
		log.Fatalf("Failed to initialize state: %v", err)
	}

	mux := http.ServeMux{}

	fileServerHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	mux.Handle("GET /app/", s.Mw.IncrementHits(fileServerHandler))

	mux.Handle("GET /admin/metrics", s.Mw.WithHits(handlers.Metrics))
	mux.Handle("POST /admin/reset", s.Mw.ResetHits(handlers.Reset(s)))

	mux.HandleFunc("GET /api/chirps/{chirpID}", handlers.ChirpLoader(s))
	mux.HandleFunc("GET /api/chirps", handlers.ChirpsLoader(s))
	mux.HandleFunc("POST /api/chirps", handlers.ChirpsAction(s))
	mux.HandleFunc("GET /api/healthz", handlers.Healthz)
	mux.HandleFunc("POST /api/users", handlers.Users(s))

	server := &http.Server{
		Handler: &mux,
		Addr:    ":8080",
	}

	fmt.Printf("Server listening on http://localhost%s\n", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
