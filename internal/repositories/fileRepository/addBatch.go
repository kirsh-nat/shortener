package filerepository

import (
	"context"
	"encoding/json"
	"os"

	"github.com/kirsh-nat/shortener.git/internal/services"
)

func (r *FileRepository) AddBatch(_ context.Context, host, userID string, data []services.BatchItem) ([]services.URLData, error) {
	var res []services.URLData

	file, err := os.OpenFile(r.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	m := make(map[string]string)
	for _, v := range data {
		short := services.MakeShortURL(v.Original)

		m[short] = v.Original

		res = append(res, services.URLData{
			ID:    v.ID,
			Short: services.MakeFullShortURL(short, host),
		})
	}

	writeData, err := json.MarshalIndent(m, "", "   ")
	if err != nil {
		return nil, err
	}
	writeData = append(writeData, '\n')

	_, err = file.Write(writeData)
	if err != nil {
		return nil, err
	}

	return res, nil
}
