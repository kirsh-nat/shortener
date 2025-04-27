package handlers

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/kirsh-nat/shortener.git/internal/app"
)

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
