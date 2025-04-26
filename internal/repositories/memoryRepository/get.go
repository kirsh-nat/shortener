package memoryRepository

import (
	"context"

	"github.com/kirsh-nat/shortener.git/internal/domain"
)

func (r *MemoryRepository) Get(_ context.Context, short string) (string, error) {
	if val, ok := r.store[short]; ok {
		return val, nil
	}

	return "", domain.ErrorURLNotFound
}
