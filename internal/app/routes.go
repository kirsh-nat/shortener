package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Routes() *chi.Mux {

	r := chi.NewRouter()

	createShortURLHandler := http.HandlerFunc(createShortURL)

	r.Post("/", http.HandlerFunc(Middleware(createShortURLHandler)))
	r.Get("/{id}", http.HandlerFunc(Middleware(http.HandlerFunc(getURL))))
	r.Post("/api/shorten", http.HandlerFunc(Middleware(http.HandlerFunc(getAPIShorten))))

	return r
}
