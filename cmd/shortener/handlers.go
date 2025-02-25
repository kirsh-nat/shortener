package main

import (
	"io"
	"net/http"

	internal "github.com/kirsh-nat/shortener.git/internal/services"
)

func createShortURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/plain")

	reqURL, _ := io.ReadAll(r.Body)
	url := string(reqURL)
	for _, v := range URLList {
		if v == url {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	shortURL := internal.NewShortURL(url)

	URLList[shortURL] = url

	_, _ = w.Write([]byte(shortURL))
	w.WriteHeader(http.StatusCreated)
}

func getURL(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	redirect, err := URLList[id]
	if !err {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)

}
