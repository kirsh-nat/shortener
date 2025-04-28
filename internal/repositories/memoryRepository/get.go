package memoryrepository

import (
	"context"

	"github.com/kirsh-nat/shortener.git/internal/domain"
)

func (r *MemoryRepository) Get(_ context.Context, short string) (string, error) {
	if val, ok := r.store[short]; ok {
		if val.Deleted {
			return val.OriginalURL, domain.NewDeletedError("URL deleted", nil)

		}
		return val.OriginalURL, nil
	}

	return "", domain.ErrorURLNotFound
}
