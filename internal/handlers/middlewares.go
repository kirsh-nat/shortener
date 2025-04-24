package handlers

import (
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/kirsh-nat/shortener.git/internal/app"
	"github.com/kirsh-nat/shortener.git/internal/models"
)

type UserKey struct{}

func Middleware(h http.Handler) http.HandlerFunc {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(&lw, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")

		cookieToken, err := r.Cookie("token")
		var user *models.User

		if err != nil || cookieToken.Value == "" {
			uuid := models.GenerateUUID()
			user, err = models.CreateUser(uuid)
			if err != nil {
				http.Error(w, "Unable to create user", http.StatusInternalServerError)
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:  "token",
				Value: user.Token,
			})
		} else {
			user, err = models.GetUser(cookieToken.Value)
			if err != nil {
				http.Error(w, "Unable to get user", http.StatusInternalServerError)
				return
			}
		}

		ctx := context.WithValue(r.Context(), UserKey{}, user)

		h.ServeHTTP(gzipWriter{ResponseWriter: &lw, Writer: gz}, r.WithContext(ctx))

		duration := time.Since(start)

		app.Sugar.Infoln(
			"type:", "request",
			"uri:", app.AppSettings.Addr+uri,
			"method:", method,
			"duration:", duration,
		)

		app.Sugar.Infoln(
			"type:", "response",
			"status:", responseData.status,
			"size:", responseData.size,
		)
	}

	return logFn
}

func GetUserFromContext(r *http.Request) (*models.User, bool) {
	user, ok := r.Context().Value(UserKey{}).(*models.User)
	return user, ok
}
