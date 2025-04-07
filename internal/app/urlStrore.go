package app

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kirsh-nat/shortener.git/internal/config"
	"github.com/kirsh-nat/shortener.git/internal/migrations"
	internal "github.com/kirsh-nat/shortener.git/internal/services"
)

const (
	//errors
	ErrURLNotFound = "URL not found"
	ErrURLExist    = "URL already exists"

	typeStorageMemory = "memory"
	typeStorageDB     = "DB"
	typeStorageFile   = "file"
)

type URLStore struct {
	mu           sync.RWMutex
	DBConnection *sql.DB
	listURL      map[string]string
	typeStorage  string
	adress       string
}

type infoURL struct {
	Decode string `json:"url"`
	Encode string `json:"short_url"`
}

type FileWriter struct {
	file   *os.File
	writer *bufio.Writer
}

type FileReader struct {
	file   *os.File
	reader *bufio.Reader
}

type urlBatchData struct {
	ID    string `json:"correlation_id"`
	Short string `json:"short_url"`
}

type DublicateError struct {
	level string
	Err   error
}

func (le *DublicateError) Error() string {
	return fmt.Sprintf("[%s] %v", le.level, le.Err)
}

func NewDublicateError(label string, err error) error {
	return &DublicateError{
		level: strings.ToUpper(label),
		Err:   err,
	}
}

func NewInfoURL() *infoURL {
	return &infoURL{}
}

// TODO: настроуки от  окружения здесь!!!!
// add typeStorage in struct : memory, file, DB !!!!!
func NewURLStore(config *config.Config) *URLStore {
	URLStore := URLStore{
		listURL:     make(map[string]string),
		typeStorage: typeStorageMemory,
		adress:      "http://" + config.Addr + "/",
	}
	if config.SetDBConnection != "" {
		URLStore.DBConnection = SetDBConnection(config.SetDBConnection)
		URLStore.typeStorage = typeStorageDB

		migrations.CreateLinkTable(URLStore.DBConnection)

	} else if config.FilePath != "" {
		reader, err := NewFileReader(config.FilePath)
		if err != nil {
			Sugar.Error(err)
			return nil
		}
		defer reader.file.Close()
		reader.ReadFile(&URLStore)
		URLStore.typeStorage = typeStorageFile
	}
	URLStore.mu.RLock()
	defer URLStore.mu.RUnlock()

	return &URLStore
}

func NewFileWriter(filename string) (*FileWriter, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &FileWriter{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

func NewFileReader(filename string) (*FileReader, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &FileReader{
		file:   file,
		reader: bufio.NewReader(file),
	}, nil
}

func (s *URLStore) Add(short, long string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	shortURL := s.adress + short
	if _, exists := s.listURL[short]; exists {
		return shortURL, NewDublicateError(s.typeStorage, errors.New(ErrURLExist))
	}

	s.listURL[short] = long
	return shortURL, nil
}

func (s *URLStore) Get(short string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	long, exists := s.listURL[short]

	if !exists {
		return "", errors.New(ErrURLNotFound)
	}

	return long, nil
}

func (s *URLStore) SaveIntoFile(short, long, fname string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	shortURL := s.adress + short
	if _, exists := s.listURL[short]; exists {
		return shortURL, NewDublicateError(s.typeStorage, errors.New(ErrURLExist))
	}

	file, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()

	m := make(map[string]string)
	m[short] = long

	data, err := json.MarshalIndent(m, "", "   ")
	if err != nil {
		return "", err
	}
	data = append(data, '\n')

	_, err = file.Write(data)
	if err != nil {
		return "", err
	}
	s.listURL[short] = long

	return shortURL, nil
}

func (r FileReader) ReadFile(url *URLStore) error {
	for {
		data, err := r.reader.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		err = json.Unmarshal(data, &url.listURL)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *URLStore) GetURLFromDBLinks(ctx context.Context, short string) (string, error) {

	row := s.DBConnection.QueryRowContext(ctx,
		"SELECT original_url from links where short_url = $1", short)
	var long sql.NullString

	err := row.Scan(&long)
	if err != nil {
		Sugar.Error(err)
		return "", err
	}
	if long.Valid {
		return long.String, nil
	}
	return "", errors.New(ErrURLNotFound)
}

func (s *URLStore) AddURLDBLinks(ctx context.Context, short, long string) (string, error) {

	_, err := s.DBConnection.ExecContext(ctx,
		"INSERT INTO links (short_url, original_url) VALUES ($1, $2)", short, long)

	s.listURL[short] = long

	shortURL := strings.TrimSpace(s.adress + short)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {

			return shortURL, NewDublicateError(s.typeStorage, err)
		}
		Sugar.Error(err)
		return "", err
	}

	return shortURL, nil
}

func (s *URLStore) InsertBatchURLsIntoDB(ctx context.Context, data []map[string]string) ([]byte, error) {
	type urlData struct {
		ID    string `json:"correlation_id"`
		Short string `json:"short_url"`
	}

	var res []urlData

	tx, err := s.DBConnection.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx,
		"INSERT INTO links (short_url, original_url) VALUES($1, $2)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	for _, v := range data {
		code := v["correlation_id"]
		original := v["original_url"]
		short := internal.MakeShortURL(original)

		_, err := stmt.ExecContext(ctx, short, original)
		if err != nil {
			return nil, err
		}
		s.Add(short, original)

		res = append(res, urlData{
			ID:    code,
			Short: s.adress + short,
		})
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	responseJSON, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	return responseJSON, nil
}

func (s *URLStore) InsertBatchURLsIntoFile(ctx context.Context, data []map[string]string, fname string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var res []urlBatchData

	file, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	m := make(map[string]string)
	for _, v := range data {
		code := v["correlation_id"]
		original := v["original_url"]
		short := internal.MakeShortURL(original)

		m[short] = original
		s.Add(short, original)

		res = append(res, urlBatchData{
			ID:    code,
			Short: s.adress + short,
		})
	}

	writeData, err := json.MarshalIndent(m, "", "   ")
	if err != nil {
		return nil, err
	}
	writeData = append(writeData, '\n')

	_, err = file.Write(writeData)
	if err != nil {
		return nil, err
	}

	responseJSON, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	return responseJSON, nil
}

func (s *URLStore) InsertBatchURLsIntoMemory(ctx context.Context, data []map[string]string) ([]byte, error) {
	var res []urlBatchData

	for _, v := range data {
		code := v["correlation_id"]
		original := v["original_url"]
		short := internal.MakeShortURL(original)
		s.Add(short, original)

		res = append(res, urlBatchData{
			ID:    code,
			Short: s.adress + short,
		})
	}

	responseJSON, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	return responseJSON, nil
}
