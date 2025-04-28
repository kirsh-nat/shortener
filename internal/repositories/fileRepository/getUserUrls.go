package filerepository

import (
	"context"
	"os"

	"github.com/kirsh-nat/shortener.git/internal/services"
)

func (r *FileRepository) GetUserURLs(ctx context.Context, userID string) ([]services.UserURLData, error) {
	file, err := os.Open(r.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data := make(map[string]services.UserURLData)
	if err := r.loadData(file, &data); err != nil {
		return nil, err
	}
	userUrls := []services.UserURLData{}

	for user, data := range data {
		if user == userID {
			userUrls = append(userUrls, data)
		}
	}

	return userUrls, nil
}
