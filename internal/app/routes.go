package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Routes() *chi.Mux {

	r := chi.NewRouter()

	createShortURLHandler := http.HandlerFunc(createShortURL)

	r.Post("/", http.HandlerFunc(WithLogging(createShortURLHandler)))
	r.Get("/{id}", http.HandlerFunc(WithLogging(http.HandlerFunc(getURL))))
	r.Get("/api/shorten", http.HandlerFunc(WithLogging(http.HandlerFunc(getApiShorten))))

	return r
}
