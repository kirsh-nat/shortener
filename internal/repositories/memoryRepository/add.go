package memoryRepository

import (
	"context"

	"github.com/kirsh-nat/shortener.git/internal/domain"
)

func (r *MemoryRepository) Add(ctx context.Context, shortURL, originalURL string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.store[shortURL]; ok {
		return domain.NewDublicateError("Memory dublicate error", nil)
	}
	r.store[shortURL] = originalURL

	return nil

}
