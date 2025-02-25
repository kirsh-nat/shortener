package main

import (
	"io"
	"net/http"

	internal "github.com/kirsh-nat/shortener.git/internal/services"
)

const localhost = "http://localhost:8080/"

func createShortURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

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

	response := localhost + shortURL

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(response))

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
	w.WriteHeader(http.StatusTemporaryRedirect)

	http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)

}
