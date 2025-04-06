package app

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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

	gzipWriter struct {
		http.ResponseWriter
		Writer io.Writer
	}
)

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

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

// TODO: вызов к структуре Store
func createShortURL(w http.ResponseWriter, r *http.Request) {
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
	response := "http://" + AppSettings.Addr + "/" + shortURL
	if Store.typeStorage == typeStorageDB {
		err = Store.AddURLDBLinks(context.Background(), shortURL, parsedURL.String())
	}
	if Store.typeStorage == typeStorageFile {
		err = Store.SaveIntoFile(shortURL, parsedURL.String(), AppSettings.FilePath)
	}
	if Store.typeStorage == typeStorageMemory {
		err = Store.Add(shortURL, parsedURL.String())
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		Sugar.Info("Can't save info in file", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(response))
}

// TODO: вынести Store на верхний уровень
func getURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Post method not allowed"))
		return
	}

	short := r.PathValue("id")

	var redirectURL string
	var err error

	if Store.typeStorage == typeStorageDB {
		redirectURL, err = Store.GetURLFromDBLinks(context.Background(), short)
	} else {
		redirectURL, err = Store.Get(short)
	}

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Location", redirectURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func getAPIShorten(w http.ResponseWriter, r *http.Request) {
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

		shortURL := internal.MakeShortURL(dataURL.URL)
		err = Store.Add(shortURL, dataURL.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			Sugar.Error(err)
			return
		}

		res := make(map[string]string, 1)
		res["result"] = "http://" + AppSettings.Resp + "/" + shortURL

		resp, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			Sugar.Error(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(resp)
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

func createBatchURLs(w http.ResponseWriter, r *http.Request) {
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
	fmt.Println("\n DATA URK \n", dataURL)

	res, err := Store.InsertBatchURLs(context.Background(), dataURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		Sugar.Error(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}
