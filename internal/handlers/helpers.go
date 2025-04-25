package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"

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

func (h *URLHandler) PingHandler(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Ping(); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *URLHandler) shortenURL(ctx context.Context, original string) (string, error) {
	shortURL := services.MakeShortURL(original)
	err := h.service.Add(ctx, shortURL, original)
	var dErr *domain.DublicateError

	if err != nil {
		if errors.As(err, &dErr) {
			return services.MakeFullShortURL(shortURL, app.AppSettings.Addr), err
		}
		return "", err
	}
	result := services.MakeFullShortURL(shortURL, app.AppSettings.Addr)

	return result, nil
}
