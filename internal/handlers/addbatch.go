package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/kirsh-nat/shortener.git/internal/app"
	"github.com/kirsh-nat/shortener.git/internal/services"
)

func (h *URLHandler) AddBatch(w http.ResponseWriter, r *http.Request) {
	if !h.checkMethod(w, r, http.MethodPost) {
		return
	}

	var dataURL []services.BatchItem

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

	var res []services.URLData

	res, err = h.service.AddBatch(r.Context(), app.AppSettings.Addr, dataURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	responseJSON, jsonErr := json.Marshal(res)
	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseJSON)
}
