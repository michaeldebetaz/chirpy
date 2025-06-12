package middlewares

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type Config struct {
	fileserverHits *atomic.Int32
}

func (c *Config) IncrementFileserverHits(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (c *Config) FileserverHits() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits := c.fileserverHits.Load()
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		body := fmt.Sprintf("Hits: %d\n", hits)
		if _, err := w.Write([]byte(body)); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	})
}

func (c *Config) ResetFileserverHits() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.fileserverHits.Store(0)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		body := fmt.Sprintf("Hits reset to 0\n")
		if _, err := w.Write([]byte(body)); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	})
}

func NewConfig() *Config {
	return &Config{
		fileserverHits: &atomic.Int32{},
	}
}
