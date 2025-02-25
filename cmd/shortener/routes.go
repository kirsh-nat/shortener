package main

import "net/http"

func routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", createShortUrl)
	mux.HandleFunc("/{id}", getUrl)

	return mux
}
