package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/kirsh-nat/shortener.git/internal/app"
	"github.com/kirsh-nat/shortener.git/internal/domain"
	"github.com/kirsh-nat/shortener.git/internal/models/user"
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

func (h *URLHandler) shortenURL(ctx context.Context, original, userID string) (string, error) {
	shortURL := services.MakeFullShortURL(services.MakeShortURL(original), app.AppSettings.Addr)
	err := h.service.Add(ctx, shortURL, original, userID)
	var dErr *domain.DublicateError

	if err != nil {
		if errors.As(err, &dErr) {
			return shortURL, err
		}
		return "", err
	}
	result := shortURL

	return result, nil
}

func (h *URLHandler) setCookieToken(w http.ResponseWriter, r *http.Request) (*user.User, bool) {
	cookieToken, err := r.Cookie("token")
	if err != nil || cookieToken == nil {
		return h.createUserAndSetCookie(w)
	}

	user, err := user.GetUser(cookieToken.Value)
	if err != nil {
		return h.createUserAndSetCookie(w)
	}

	return user, true
}

func (h *URLHandler) createUserAndSetCookie(w http.ResponseWriter) (*user.User, bool) {
	user, err := user.CreateUser()
	if err != nil {
		return nil, false
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: user.Token,
		Path:  "/",
	})
	return user, true
}

func (h *URLHandler) getCookieToken(w http.ResponseWriter, r *http.Request) (*user.User, bool) {
	cookieToken, err := r.Cookie("token")
	if err != nil || cookieToken == nil {
		return h.unauthorizedResponse(w)
	}

	user, err := user.GetUser(cookieToken.Value)
	if err != nil {
		return h.unauthorizedResponse(w)
	}

	return user, true
}

func (h *URLHandler) unauthorizedResponse(w http.ResponseWriter) (*user.User, bool) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("Unauthorized"))
	return nil, false
}
