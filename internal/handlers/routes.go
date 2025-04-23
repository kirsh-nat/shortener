package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	//"github.com/kirsh-nat/shortener.git/internal/handlers"
	//"github.com/kirsh-nat/shortener.git/internal/handlers"
	//"github.com/kirsh-nat/shortener.git/internal/handlers"
)

func Routes(handler *URLHandler) *chi.Mux {

	r := chi.NewRouter()

	r.Post("/", http.HandlerFunc(Middleware(http.HandlerFunc(handler.Add))))
	r.Get("/{id}", handler.Get)
	r.Post("/api/shorten", http.HandlerFunc(Middleware(http.HandlerFunc(handler.GetAPIShorten))))
	r.Get("/ping", http.HandlerFunc(Middleware(http.HandlerFunc(handler.PingHandler))))
	r.Post("/api/shorten/batch", http.HandlerFunc(Middleware(http.HandlerFunc(handler.AddBatch))))
	r.Get("/api/user/urls", http.HandlerFunc(Middleware(http.HandlerFunc(handler.GetUserURLs))))
	r.Delete("/api/user/urls", http.HandlerFunc(Middleware(http.HandlerFunc(handler.DeleteUserURLs))))

	return r
}
