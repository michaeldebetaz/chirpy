package middlewares

import (
	"context"
	"net/http"
	"sync/atomic"
)

type contextKey string

const HITS_KEY contextKey = "hits"

type Middleware struct {
	Hits *atomic.Int32
}

func (m *Middleware) IncrementHits(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.Hits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) WithHits(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits := m.Hits.Load()
		ctx := context.WithValue(r.Context(), HITS_KEY, hits)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) ResetHits(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.Hits.Store(0)

		hits := m.Hits.Load()
		ctx := context.WithValue(r.Context(), HITS_KEY, hits)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func New() *Middleware {
	return &Middleware{
		Hits: &atomic.Int32{},
	}
}
