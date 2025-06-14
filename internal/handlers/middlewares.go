package handlers

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/kirsh-nat/shortener.git/internal/app"
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

		// cookieToken, err := r.Cookie("token")
		// var user *models.User

		// if err != nil || cookieToken == nil || cookieToken.Value == "" {
		// 	uuid := models.GenerateUUID()
		// 	app.Sugar.Info("created user UUID: ", uuid)
		// 	user, err = models.CreateUser(uuid)
		// 	if err != nil {
		// 		http.Error(w, "Unable to create user", http.StatusInternalServerError)
		// 		return
		// 	}
		// 	http.SetCookie(w, &http.Cookie{
		// 		Name:  "token",
		// 		Value: user.Token,
		// 	})
		// } else {
		// 	user, err = models.GetUser(cookieToken.Value)
		// 	if err != nil {
		// 		http.Error(w, "Unable to get user", http.StatusInternalServerError)
		// 		return
		// 	}
		// }

		//ctx := context.WithValue(r.Context(), UserKey{}, user)

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(&lw, r)
			//h.ServeHTTP(&lw, r.WithContext(ctx))
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")

		//h.ServeHTTP(gzipWriter{ResponseWriter: &lw, Writer: gz}, r.WithContext(ctx))
		h.ServeHTTP(gzipWriter{ResponseWriter: &lw, Writer: gz}, r)

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
