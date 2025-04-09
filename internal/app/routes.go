package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Routes(s *URLStore) *chi.Mux {

	r := chi.NewRouter()

	createShortURLHandler := http.HandlerFunc(s.createShortURL)
	r.Post("/", http.HandlerFunc(Middleware(createShortURLHandler)))
	r.Get("/{id}", http.HandlerFunc(Middleware(http.HandlerFunc(s.getURL))))
	r.Post("/api/shorten", http.HandlerFunc(Middleware(http.HandlerFunc(s.getAPIShorten))))
	r.Post("/api/shorten/batch", http.HandlerFunc(Middleware(http.HandlerFunc(s.createBatchURLs))))
	r.Get("/ping", http.HandlerFunc(Middleware(http.HandlerFunc(pingHandler))))

	return r
}
