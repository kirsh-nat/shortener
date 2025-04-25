package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kirsh-nat/shortener.git/internal/app"
	"github.com/kirsh-nat/shortener.git/internal/domain"
	"github.com/kirsh-nat/shortener.git/internal/services"
)

func (h *URLHandler) GetAPIShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
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

		w.Header().Set("Content-Type", "application/json")
		shortURL := services.MakeShortURL(dataURL.URL)
		err = h.service.Add(r.Context(), shortURL, dataURL.URL)
		var response []byte
		result := "http://" + app.AppSettings.Addr + "/" + shortURL
		var dErr *domain.DublicateError

		res := make(map[string]string, 1)
		res["result"] = result

		var jsonErr error
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

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

}
