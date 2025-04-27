package handlers

import (
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/kirsh-nat/shortener.git/internal/domain"
)

func (h *URLHandler) Add(w http.ResponseWriter, r *http.Request) {
	if !h.checkMethod(w, r, http.MethodPost) {
		return
	}
	var body io.Reader = r.Body
	if r.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Can't create gzip reader"))
			return
		}
		defer gz.Close()
		body = gz
	}

	reqURL, err := io.ReadAll(body)
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

	shortURL, err := h.shortenURL(r.Context(), parsedURL.String())
	var dErr *domain.DublicateError
	var response string
	if errors.As(err, &dErr) {
		response = shortURL
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(response))
		return

	} else if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if response == "" {
		response = shortURL
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(response))
}
