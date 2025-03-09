package app

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

	reqURL, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	url := string(reqURL)
	shortURL, err := internal.AddURL(url, &listURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := "http://" + AppSettings.Addr + "/" + shortURL

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(response))

}

func getURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	redirectURL, err := internal.GetURL(id, listURL)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Location", redirectURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
