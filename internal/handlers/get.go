package handlers

import (
	"net/http"

	"github.com/kirsh-nat/shortener.git/internal/app"
	"github.com/kirsh-nat/shortener.git/internal/services"
)

func (h *URLHandler) Get(w http.ResponseWriter, r *http.Request) {
	if !h.checkMethod(w, r, http.MethodGet) {
		return
	}

	short := r.PathValue("id")

	var redirectURL string
	var err error

	redirectURL, err = h.service.Get(r.Context(), services.MakeFullShortURL(short, app.AppSettings.Addr))

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Location", redirectURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
