package app

import (
	"io"
	"net/http"
	"net/url"

	internal "github.com/kirsh-nat/shortener.git/internal/services"
)

func createShortURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))

		return
	}

	reqURL, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Can't read request body"))

		return
	}

	parsedURL, err := url.ParseRequestURI(string(reqURL))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	shortURL := internal.MakeShortURL(parsedURL.String())
	err = Store.Add(shortURL, parsedURL.String())

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	response := "http://" + AppSettings.Addr + "/" + shortURL

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(response))

}

func getURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(" Post method not allowed"))

		return
	}

	short := r.PathValue("id")
	redirectURL, err := Store.Get(short)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))

		return
	}

	w.Header().Set("Location", redirectURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
