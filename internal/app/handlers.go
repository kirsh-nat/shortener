package app

import (
	"io"
	"net/http"
	"net/url"
	"time"

	internal "github.com/kirsh-nat/shortener.git/internal/services"
)

type (
	// структура для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}

	// реализация http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

func WithLogging(h http.Handler) http.HandlerFunc {
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

		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		Sugar.Infoln(
			"type:", "request",
			"uri:", AppSettings.Addr+uri,
			"method:", method,
			"duration:", duration,
		)

		Sugar.Infoln(
			"type:", "response",
			"status:", responseData.status,
			"size:", responseData.size,
		)
	}

	return logFn
}

func createShortURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))

		return
	}

	reqURL, err := io.ReadAll(r.Body)
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

	shortURL := internal.MakeShortURL(parsedURL.String())
	err = Store.Add(shortURL, parsedURL.String())

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	response := "http://" + AppSettings.Addr + "/" + shortURL

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(response))

}

func getURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(" Post method not allowed"))

		return
	}

	short := r.PathValue("id")
	redirectURL, err := Store.Get(short)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))

		return
	}

	w.Header().Set("Location", redirectURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
