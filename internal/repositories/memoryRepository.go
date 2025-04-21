package repositories

import (
	"context"
	"encoding/json"

	"github.com/kirsh-nat/shortener.git/internal/domain"
	"github.com/kirsh-nat/shortener.git/internal/models"
	"github.com/kirsh-nat/shortener.git/internal/services"
)

type MemoryRepository struct {
	store map[string]string
}

func NewMemoryRepository() models.URLRepository {
	return &MemoryRepository{store: make(map[string]string)}
}

func (r *MemoryRepository) Add(ctx context.Context, shortURL, originalURL string) error {
	if _, ok := r.store[shortURL]; ok {
		return domain.NewDublicateError("Memory dublicate error", nil)
	}
	r.store[shortURL] = originalURL

	return nil

}

func (r *MemoryRepository) Get(short string) (string, error) {
	if val, ok := r.store[short]; ok {
		return val, nil
	}

	return "", domain.ErrorURLNotFound
}

func (r *MemoryRepository) Ping() error {
	return nil
}

func (r *MemoryRepository) AddBatch(host string, data []map[string]string) ([]byte, error) {
	type urlData struct {
		ID    string `json:"correlation_id"`
		Short string `json:"short_url"`
	}

	var res []urlData

	for _, v := range data {
		code := v["correlation_id"]
		original := v["original_url"]
		short := services.MakeShortURL(original)

		err := r.Add(context.Background(), short, original)
		if err != nil {
			return nil, err
		}

		res = append(res, urlData{
			ID:    code,
			Short: "http://" + host + "/" + short,
		})
	}

	responseJSON, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	return responseJSON, nil
}
