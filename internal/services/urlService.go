package services

import (
	"context"

	"github.com/kirsh-nat/shortener.git/internal/models"
)

type URLService struct {
	repo     models.URLRepository
	userURLs map[string][]string
}

func NewURLService(repo models.URLRepository) *URLService {
	return &URLService{repo: repo, userURLs: make(map[string][]string)}
}

func (s *URLService) Add(ctx context.Context, shortURL, originalURL string) error {
	err := s.repo.Add(ctx, shortURL, originalURL)
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

func (s *URLService) AddUserURL(userId, short string) {
	if _, ok := s.userURLs[userId]; !ok {
		s.userURLs[userId] = make([]string, 0)
	}

	s.userURLs[userId] = append(s.userURLs[userId], short)
}

func (s *URLService) GetUserURLs(userId string) []string {

	if _, ok := s.userURLs[userId]; !ok {
		return []string{}
	}

	return s.userURLs[userId]
}
