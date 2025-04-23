package models

import "context"

type URLRepository interface {
	Add(ctx context.Context, shortURL, originalURL, userID string) error
	Get(short string) (string, error)
	Ping() error
	AddBatch(host string, data []map[string]string) ([]byte, error)
	DeleteBatch(shortURLs []string, userID string)
}
