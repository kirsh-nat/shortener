package services

import (
	"context"

	"github.com/kirsh-nat/shortener.git/internal/models"
)

type URLService struct {
	repo models.URLRepository // интерфейс для доступа к данным
}

func NewURLService(repo models.URLRepository) *URLService {
	return &URLService{repo: repo}
}

// Реализация метода Add
func (s *URLService) Add(ctx context.Context, shortURL, originalURL string) error {
	// Логика добавления URL в репозиторий
	err := s.repo.Add(shortURL, originalURL)
	if err != nil {
		return err
	}
	return nil
}

// Реализация метода Get
func (s *URLService) Get(ctx context.Context, short string) (string, error) {
	// Логика получения оригинального URL из репозитория
	longURL, err := s.repo.Get(short)
	if err != nil {
		return "", err
	}
	return longURL, nil
}

func (s *URLService) Ping() error {
	return s.repo.Ping()
}

func (s *URLService) AddBatch(data []map[string]string) ([]byte, error) {
	return s.repo.AddBatch(data)
}
