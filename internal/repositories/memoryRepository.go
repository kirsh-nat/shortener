package repositories

import (
	"context"
	"sync"

	"github.com/kirsh-nat/shortener.git/internal/domain"
	"github.com/kirsh-nat/shortener.git/internal/services"
)

type MemoryRepository struct {
	mu    sync.RWMutex
	store map[string]string
}

func NewMemoryRepository() services.URLRepository {
	return &MemoryRepository{store: make(map[string]string)}
}

func (r *MemoryRepository) Add(ctx context.Context, shortURL, originalURL string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.store[shortURL]; ok {
		return domain.NewDublicateError("Memory dublicate error", nil)
	}
	r.store[shortURL] = originalURL

	return nil

}

func (r *MemoryRepository) Get(_ context.Context, short string) (string, error) {
	if val, ok := r.store[short]; ok {
		return val, nil
	}

	return "", domain.ErrorURLNotFound
}

func (r *MemoryRepository) Ping() error {
	return nil
}

func (r *MemoryRepository) AddBatch(context context.Context, host string, data []services.BatchItem) ([]services.URLData, error) {
	var res []services.URLData

	for _, v := range data {
		short := services.MakeShortURL(v.Original)

		err := r.Add(context, short, v.Original)
		if err != nil {
			return nil, err
		}

		res = append(res, services.URLData{
			ID:    v.ID,
			Short: "http://" + host + "/" + short,
		})
	}

	return res, nil
}
