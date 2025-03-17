package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Routes() *chi.Mux {

	r := chi.NewRouter()

	//TODO: сделать элегантнее!!!
	createShortURLHandler := http.HandlerFunc(createShortURL)

	r.Post("/", http.HandlerFunc(WithLogging(createShortURLHandler))) //createShortURL
	r.Get("/{id}", http.HandlerFunc(WithLogging(http.HandlerFunc(getURL))))
	r.Post("/api/shorten", http.HandlerFunc(WithLogging(http.HandlerFunc(getApiURL))))

	return r
}
