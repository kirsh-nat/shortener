package models

type URLRepository interface {
	Add(shortURL, originalURL string) error
	Get(short string) (string, error)
	Ping() error
	AddBatch(host string, data []map[string]string) ([]byte, error)
}
