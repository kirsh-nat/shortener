package main

import (
	"github.com/go-chi/chi/v5"
)

func routes() *chi.Mux { // *http.ServeMux
	r := chi.NewRouter()
	r.Post("/", createShortURL)
	r.Get("/{id}", getURL)

	return r
}
