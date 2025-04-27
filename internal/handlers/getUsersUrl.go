package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kirsh-nat/shortener.git/internal/app"
)

func (h *URLHandler) GetUserURLs(w http.ResponseWriter, r *http.Request) {
	if !h.checkMethod(w, r, http.MethodGet) {
		return
	}

	user, ok := h.setCookieToken(w, r)
	if !ok {
		return
	}

	shortUrls, err := h.service.GetUserURLs(r.Context(), user.UUID)

	app.Sugar.Debug("Get user urls: ", user.UUID, " shortUrls ", shortUrls)
	if err != nil {
		app.Sugar.Errorw(err.Error(), "event", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Can't get user urls"))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if len(shortUrls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	resp, jsonErr := json.Marshal(shortUrls)

	if jsonErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	app.Sugar.Debug("Json user urls: ", user.UUID, " shortUrls ", resp)

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
