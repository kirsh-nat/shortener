package handlers

import "net/http"

func (h *URLHandler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Post method not allowed"))
		return
	}

	short := r.PathValue("id")

	var redirectURL string
	var err error

	redirectURL, err = h.service.Get(r.Context(), short)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Location", redirectURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
