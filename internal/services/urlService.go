package services

import (
	"context"

	"github.com/kirsh-nat/shortener.git/internal/models"
)

type URLService struct {
	repo models.URLRepository
}

func NewURLService(repo models.URLRepository) *URLService {
	return &URLService{repo: repo}
}

func (s *URLService) Add(ctx context.Context, shortURL, originalURL, userID string) error {
	err := s.repo.Add(ctx, shortURL, originalURL, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *URLService) Get(ctx context.Context, short string) (string, error) {
	longURL, err := s.repo.Get(short)
	if err != nil {
		return "", err
	}
	return longURL, nil
}

func (s *URLService) Ping() error {
	return s.repo.Ping()
}

func (s *URLService) AddBatch(host string, data []map[string]string) ([]byte, error) {
	return s.repo.AddBatch(host, data)
}

func (s *URLService) DeleteBatch(shortURLs []string, userID string) {
	s.repo.DeleteBatch(shortURLs, userID)
}

func (s *URLService) AddUserURL(userID, short string) {
	s.repo.AddUserURL(userID, short)
}

func (s *URLService) GetUserURLs(userID string) ([]string, error) {
	return s.repo.GetUserURLs(userID)
}
