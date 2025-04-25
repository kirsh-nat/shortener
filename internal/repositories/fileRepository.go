package repositories

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/kirsh-nat/shortener.git/internal/domain"
	"github.com/kirsh-nat/shortener.git/internal/services"
)

type FileRepository struct {
	//db       *sql.DB
	filePath string
}

type FileReader struct {
	file   *os.File
	reader *bufio.Reader
}

func NewFileRepository(filePath string) services.URLRepository {
	return &FileRepository{filePath: filePath}
}

func (r *FileRepository) Add(ctx context.Context, shortURL, originalURL string) error {
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

func (r *FileRepository) loadData(file *os.File, data *map[string]string) error {
	file.Seek(0, 0)

	var temp map[string]string
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

func (r *FileRepository) AddBatch(_ context.Context, host string, data []services.BatchItem) ([]services.URLData, error) {
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
			Short: "http://" + host + "/" + short,
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

func (r FileReader) readFile(res *map[string]string) error {
	b := make(map[string]string)
	for {
		data, err := r.reader.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}

		err = json.Unmarshal(data, &res)
		if err != nil {
			return err
		}

		for k, v := range b {
			(*res)[k] = v
		}
		b = make(map[string]string)
	}

	return nil
}

func newFileReader(filename string) (*FileReader, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &FileReader{
		file:   file,
		reader: bufio.NewReader(file),
	}, nil
}
