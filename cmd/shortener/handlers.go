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

	reqURL, _ := io.ReadAll(r.Body)
	url := string(reqURL)
	//TODO: check by regular
	for _, v := range URLList {
		if v == url {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	shortURL := internal.NewShortURL(url)
	URLList[shortURL] = url
	response := "http://" + conf.Resp + "/" + shortURL

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(response))

}

func getURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	redirectURL, err := URLList[id]
	if !err {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Location", redirectURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
