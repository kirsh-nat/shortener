package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/kirsh-nat/shortener.git/internal/app"
	"github.com/kirsh-nat/shortener.git/internal/services"
)

func (h *URLHandler) DeleteUserURLs(w http.ResponseWriter, r *http.Request) {
	check := h.checkMethod(w, r, http.MethodDelete)
	if !check {
		return
	}

	user, ok := h.getCookieToken(w, r)
	if !ok {
		return
	}

	var dataURL []string
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

	var deleteList []string
	for _, url := range dataURL {
		deleteList = append(deleteList, services.MakeFullShortURL(url, app.AppSettings.Addr))
	}

	go h.service.DeleteBatch(deleteList, user.UUID)

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("OK"))
}
