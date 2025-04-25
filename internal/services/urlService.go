package services

import (
	"context"
)

type BatchItem struct {
	ID       string `json:"correlation_id"`
	Original string `json:"original_url"`
}

type URLRepository interface {
	Add(ctx context.Context, shortURL, originalURL string) error
	Get(context context.Context, short string) (string, error)
	Ping() error
	AddBatch(context context.Context, host string, data []BatchItem) ([]URLData, error)
}

type URLData struct {
	ID    string `json:"correlation_id"`
	Short string `json:"short_url"`
}

type URLService struct {
	repo URLRepository
}

func NewURLService(repo URLRepository) *URLService {
	return &URLService{repo: repo}
}

func (s *URLService) Add(ctx context.Context, shortURL, originalURL string) error {
	err := s.repo.Add(ctx, shortURL, originalURL)
	if err != nil {
		return err
	}
	return nil
}

func (s *URLService) Get(ctx context.Context, short string) (string, error) {
	longURL, err := s.repo.Get(ctx, short)
	if err != nil {
		return "", err
	}
	return longURL, nil
}

func (s *URLService) Ping() error {
	return s.repo.Ping()
}

func (s *URLService) AddBatch(context context.Context, host string, data []BatchItem) ([]URLData, error) {
	return s.repo.AddBatch(context, host, data)
}
