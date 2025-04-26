package filerepository

import (
	"context"
	"os"

	"github.com/kirsh-nat/shortener.git/internal/domain"
)

func (r *FileRepository) Get(_ context.Context, short string) (string, error) {
	file, err := os.Open(r.filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	data := make(map[string]string)
	if err := r.loadData(file, &data); err != nil {
		return "", err
	}

	if val, ok := data[short]; ok {
		return val, nil
	}

	return "", domain.ErrorURLNotFound
}
