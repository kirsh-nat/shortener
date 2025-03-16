package app

import (
	"errors"
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

func NewURLStore() *URLStore {
	return &URLStore{
		listURL: make(map[string]string),
	}
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
