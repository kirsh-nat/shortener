package filerepository

import (
	"context"

	"github.com/kirsh-nat/shortener.git/internal/services"
)

func (r *FileRepository) GetUserURLs(ctx context.Context, userID string) ([]services.UserURLData, error) {

	userUrls := []services.UserURLData{}

	return userUrls, nil
}
