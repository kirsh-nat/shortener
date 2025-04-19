package handlers

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/kirsh-nat/shortener.git/internal/app"
	"github.com/kirsh-nat/shortener.git/internal/domain"
	"github.com/kirsh-nat/shortener.git/internal/services"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}

	gzipWriter struct {
		http.ResponseWriter
		Writer io.Writer
	}
)

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

type URLHandler struct {
	service services.URLService
}

func NewURLHandler(service *services.URLService) *URLHandler {
	return &URLHandler{service: *service}
}

func (h *URLHandler) Add(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
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

	//TODO: передать хост нормально
	shortURL := services.MakeShortURL(parsedURL.String())
	w.Header().Set("Content-Type", "application/json")
	err = h.service.Add(context.Background(), shortURL, parsedURL.String())
	var dErr *domain.DublicateError
	var response string
	if errors.As(err, &dErr) {
		response = "http://" + app.AppSettings.Addr + "/" + shortURL
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(response))
		return

	} else if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if response == "" {
		response = "http://" + app.AppSettings.Addr + "/" + shortURL
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(response))
}

func (h *URLHandler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Post method not allowed"))
		return
	}

	short := r.PathValue("id")

	var redirectURL string
	var err error

	redirectURL, err = h.service.Get(context.Background(), short)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Location", redirectURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

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
		err = h.service.Add(context.Background(), shortURL, dataURL.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var response []byte
		result := "http://" + app.AppSettings.Addr + "/" + shortURL

		res := make(map[string]string, 1)
		res["result"] = result

		var jsonErr error
		response, jsonErr = json.Marshal(res)
		if jsonErr != nil {
			http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
			return
		}

		var dErr *domain.DublicateError
		if errors.As(err, &dErr) {
			w.WriteHeader(http.StatusConflict)
			w.Write(response)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(response)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

}

func (h *URLHandler) PingHandler(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Ping(); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *URLHandler) AddBatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}

	var dataURL []map[string]string

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

	var res []byte

	res, err = h.service.AddBatch(dataURL)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}
