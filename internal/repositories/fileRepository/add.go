package filerepository

import (
	"context"
	"encoding/json"
	"os"
)

// TODO: make user url store
func (r *FileRepository) Add(ctx context.Context, shortURL, originalURL, _ string) error {
	file, err := os.OpenFile(r.filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	data := make(map[string]string)
	if err := r.loadData(file, &data); err != nil {
		return err
	}

	data[shortURL] = originalURL

	tempFilePath := r.filePath + ".tmp"
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	newData, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		return err
	}
	newData = append(newData, '\n')

	if _, err := tempFile.Write(newData); err != nil {
		return err
	}

	if err := tempFile.Close(); err != nil {
		return err
	}

	if err := os.Rename(tempFilePath, r.filePath); err != nil {
		return err
	}

	return nil
}
