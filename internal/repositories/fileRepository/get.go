package filerepository

import (
	"context"
	"os"

	"github.com/kirsh-nat/shortener.git/internal/domain"
	"github.com/kirsh-nat/shortener.git/internal/services"
)

func (r *FileRepository) Get(_ context.Context, short string) (string, error) {
	file, err := os.Open(r.filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	data := make(map[string]services.UserURLData)
	if err := r.loadData(file, &data); err != nil {
		return "", err
	}

	for _, v := range data {
		if v.Short == short {
			return v.Original, nil
		}
	}

	return "", domain.ErrorURLNotFound
}
