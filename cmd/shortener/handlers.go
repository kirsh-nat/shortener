package main

import (
	"io"
	"net/http"

	internal "github.com/kirsh-nat/shortener.git/internal/services"
)

func createShortUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/plain")

	reqUrl, _ := io.ReadAll(r.Body)
	url := string(reqUrl)
	for _, v := range UrlList {
		if v == url {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	shortUrl := internal.NewShortUrl(url)

	UrlList[shortUrl] = url

	_, _ = w.Write([]byte(shortUrl))
	w.WriteHeader(http.StatusCreated)
}

func getUrl(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	redirect, err := UrlList[id]
	if !err {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)

}
