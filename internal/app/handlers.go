package app

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	internal "github.com/kirsh-nat/shortener.git/internal/services"
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

func (s *URLStore) createShortURL(w http.ResponseWriter, r *http.Request) {
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

	shortURL := internal.MakeShortURL(parsedURL.String())
	w.Header().Set("Content-Type", "application/json")
	response, err := s.Add(shortURL, parsedURL.String())
	var dErr *DublicateError
	if errors.As(err, &dErr) {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(response))
		return

	} else if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(response))
}

func (s *URLStore) getURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Post method not allowed"))
		return
	}

	short := r.PathValue("id")

	var redirectURL string
	var err error

	if s.typeStorage == typeStorageDB {
		redirectURL, err = s.GetURLFromDBLinks(context.Background(), short)
	} else {
		redirectURL, err = s.Get(short)
	}

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Location", redirectURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (s *URLStore) getAPIShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var dataURL struct {
			URL string `json:"url"`
		}

		var buf bytes.Buffer
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			Sugar.Error(err)
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &dataURL); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			Sugar.Error(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		shortURL := internal.MakeShortURL(dataURL.URL)
		result, err := s.Add(shortURL, dataURL.URL)
		var response []byte
		if result != "" {
			res := make(map[string]string, 1)
			res["result"] = result

			var jsonErr error
			response, jsonErr = json.Marshal(res)
			if jsonErr != nil {
				http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
				Sugar.Error(jsonErr)
				return
			}
		}

		var dErr *DublicateError
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
		Sugar.Infoln("request error method: %v not allowed", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	if err := DB.Ping(); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		Sugar.Error("Database connection error:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *URLStore) createBatchURLs(w http.ResponseWriter, r *http.Request) {
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
		Sugar.Error(err)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &dataURL); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		Sugar.Error(err)
		return
	}

	var res []byte

	if s.typeStorage == typeStorageDB {
		res, err = s.InsertBatchURLsIntoDB(dataURL)
	}
	if s.typeStorage == typeStorageFile {
		res, err = s.InsertBatchURLsIntoFile(dataURL, AppSettings.FilePath)
	}
	if s.typeStorage == typeStorageMemory {
		res, err = s.InsertBatchURLsIntoMemory(dataURL)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		Sugar.Error(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}
