package app

import (
	"github.com/go-chi/chi/v5"
)

func Routes() *chi.Mux {

	r := chi.NewRouter()
	r.Post("/", createShortURL)
	r.Get("/{id}", getURL)

	return r
}
