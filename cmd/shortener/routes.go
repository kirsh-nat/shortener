package main

import "net/http"

func routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", createShortURL)
	mux.HandleFunc("/{id}", getURL)

	return mux
}
