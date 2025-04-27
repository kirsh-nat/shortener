package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kirsh-nat/shortener.git/internal/app"
	"github.com/kirsh-nat/shortener.git/internal/domain"
)

func (h *URLHandler) GetAPIShorten(w http.ResponseWriter, r *http.Request) {
	if !h.checkMethod(w, r, http.MethodPost) {
		return
	}

	var dataURL struct {
		URL string `json:"url"`
	}

	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &dataURL); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, ok := h.setCookieToken(w, r)
	if !ok {
		return
	}

	app.Sugar.Debug("ADD SHORTEN user urls: ", user.UUID, " parsedURL ", dataURL.URL)

	result, err := h.shortenURL(r.Context(), dataURL.URL, user.UUID)

	var dErr *domain.DublicateError
	res := make(map[string]string, 1)
	res["result"] = result

	var jsonErr error
	var response []byte
	w.Header().Set("Content-Type", "application/json")

	response, jsonErr = json.Marshal(res)
	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
		return
	}

	if err != nil {
		if errors.As(err, &dErr) {
			w.WriteHeader(http.StatusConflict)
			w.Write(response)
			return
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
