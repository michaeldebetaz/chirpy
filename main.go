package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/michaeldebetaz/chirpy/internal/config"
	"github.com/michaeldebetaz/chirpy/internal/middlewares"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("Failed to initialize state: %v", err)
	}

	mux := http.ServeMux{}

	fileServerHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	mux.Handle("GET /app/", cfg.Mw.IncrementHits(fileServerHandler))

	mux.Handle("GET /admin/metrics", cfg.Mw.WithHits(metrics))
	mux.Handle("POST /admin/reset", cfg.Mw.ResetHits(cfg.ResetAction))

	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.ChirpLoader)
	mux.HandleFunc("GET /api/chirps", cfg.ChirpsLoader)
	mux.HandleFunc("POST /api/chirps", cfg.ChirpsAction)
	mux.HandleFunc("GET /api/healthz", healthz)
	mux.HandleFunc("POST /api/login", cfg.LoginAction)
	mux.HandleFunc("POST /api/refresh", cfg.RefreshAction)
	mux.HandleFunc("POST /api/revoke", cfg.RevokeAction)
	mux.HandleFunc("POST /api/users", cfg.UsersAction)

	server := &http.Server{
		Handler: &mux,
		Addr:    ":8080",
	}

	fmt.Printf("Server listening on http://localhost%s\n", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	body := fmt.Sprintf("%s", http.StatusText(http.StatusOK))
	w.Write([]byte(body))
}

func metrics(w http.ResponseWriter, r *http.Request) {
	hits, ok := r.Context().Value(middlewares.HITS_KEY).(int32)
	if !ok {
		http.Error(w, "Failed to retrieve hits from context", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	body := fmt.Sprintf(`
<html> 
	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p> 
	</body>
</html>`, hits)
	if _, err := w.Write([]byte(body)); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}
