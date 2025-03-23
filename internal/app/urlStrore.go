package app

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

const (
	ErrURLNotFound = "URL not found"
	ErrURLExist    = "URL already exists"
)

type URLStore struct {
	mu      sync.RWMutex
	listURL map[string]string
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

func NewInfoURL() *infoURL {
	return &infoURL{}
}

func NewURLStore(fname string) *URLStore {
	URLStore := URLStore{
		listURL: make(map[string]string),
	}
	reader, err := NewFileReader(fname)
	if err != nil {
		Sugar.Error(err)
		return nil
	}

	reader.ReadFile(&URLStore)
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
	fmt.Println(filename, 1111111111111111)
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &FileReader{
		file:   file,
		reader: bufio.NewReader(file),
	}, nil
}

func (s *URLStore) Add(short, long string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.listURL[short]; exists {
		return errors.New(ErrURLExist)
	}

	s.listURL[short] = long
	return nil
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

func (s *URLStore) SaveIntoFile(short, long, fname string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.listURL[short]; exists {
		return errors.New(ErrURLExist)
	}

	file, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	m := make(map[string]string)
	m[short] = long

	data, err := json.MarshalIndent(m, "", "   ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	_, err = file.Write(data)
	s.listURL[short] = long

	if err != nil {
		return err
	}

	return nil
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
