package repositories

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/kirsh-nat/shortener.git/internal/domain"
	"github.com/kirsh-nat/shortener.git/internal/models"
	"github.com/kirsh-nat/shortener.git/internal/services"
)

type MemoryRepository struct {
	mu       sync.RWMutex
	store    map[string]string
	userURLs map[string][]string
}

func NewMemoryRepository() models.URLRepository {
	return &MemoryRepository{store: make(map[string]string), userURLs: make(map[string][]string)}
}

func (r *MemoryRepository) Add(ctx context.Context, shortURL, originalURL, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

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

		err := r.Add(context.Background(), short, original, "")
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

func (r *MemoryRepository) DeleteBatch(data []string, userID string) {
}

func (r *MemoryRepository) AddUserURL(userID, short string) {
	if _, ok := r.userURLs[userID]; !ok {
		r.userURLs[userID] = make([]string, 0)
	}

	r.userURLs[userID] = append(r.userURLs[userID], short)
}

func (r *MemoryRepository) GetUserURLs(userID string) ([]string, error) {
	if _, ok := r.userURLs[userID]; !ok {
		return []string{}, nil
	}

	return r.userURLs[userID], nil
}
