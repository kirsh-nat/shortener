package memoryrepository

import (
	"context"

	"github.com/kirsh-nat/shortener.git/internal/services"
)

func (r *MemoryRepository) GetUserURLs(ctx context.Context, userID string) ([]services.UserURLData, error) {
	userUrls := []services.UserURLData{}
	for short, urlData := range r.store {
		if urlData.UserID == userID {
			userUrls = append(userUrls, services.UserURLData{
				Short:    short,
				Original: urlData.OriginalURL,
			})
		}
	}

	return userUrls, nil
}
