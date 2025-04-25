package models

import "context"

type URLRepository interface {
	Add(ctx context.Context, shortURL, originalURL string) error
	Get(context context.Context, short string) (string, error)
	Ping() error
	AddBatch(context context.Context, host string, data []map[string]string) ([]byte, error)
}
