package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Routes() *chi.Mux { //TODO: add AppConfig as param

	r := chi.NewRouter()

	createShortURLHandler := http.HandlerFunc(createShortURL)
	r.Post("/", http.HandlerFunc(Middleware(createShortURLHandler)))
	r.Get("/{id}", http.HandlerFunc(Middleware(http.HandlerFunc(getURL))))
	r.Post("/api/shorten", http.HandlerFunc(Middleware(http.HandlerFunc(getAPIShorten))))
	r.Post("/api/shorten/batch", http.HandlerFunc(Middleware(http.HandlerFunc(createBatchURLs))))
	r.Get("/ping", http.HandlerFunc(Middleware(http.HandlerFunc(pingHandler))))

	return r
}
