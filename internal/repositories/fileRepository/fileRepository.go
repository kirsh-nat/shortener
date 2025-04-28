package filerepository

import (
	"encoding/json"
	"io"
	"os"

	"github.com/kirsh-nat/shortener.git/internal/services"
)

type FileRepository struct {
	filePath string
}

func NewFileRepository(filePath string) services.URLRepository {
	return &FileRepository{filePath: filePath}
}

func (r *FileRepository) loadData(file *os.File, data *map[string]services.UserURLData) error {
	file.Seek(0, 0)

	var temp map[string]services.UserURLData
	if err := json.NewDecoder(file).Decode(&temp); err != nil && !os.IsNotExist(err) {
		if err != io.EOF {
			return err
		}
	}

	for k, v := range temp {
		(*data)[k] = v
	}

	return nil
}

func (r *FileRepository) Ping() error {
	return nil
}
