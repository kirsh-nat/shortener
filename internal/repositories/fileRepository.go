package repositories

import (
	"bufio"
	"encoding/json"
	"io"
	"os"

	"github.com/kirsh-nat/shortener.git/internal/domain"
	"github.com/kirsh-nat/shortener.git/internal/models"
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

type urlBatchData struct {
	ID    string `json:"correlation_id"`
	Short string `json:"short_url"`
}

func NewFileRepository(filePath string) models.URLRepository {
	return &FileRepository{filePath: filePath}
}

func (r *FileRepository) Add(shortURL, originalURL string) error {
	file, err := os.OpenFile(r.filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	data := make(map[string]string)
	if err := r.loadData(file, &data); err != nil {
		return err
	}

	data[shortURL] = originalURL

	if err := file.Truncate(0); err != nil {
		return err
	}
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	newData, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		return err
	}
	newData = append(newData, '\n')

	if _, err := file.Write(newData); err != nil {
		return err
	}

	return nil
}

func (r *FileRepository) Get(short string) (string, error) {
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
		return err
	}

	for k, v := range temp {
		(*data)[k] = v
	}
	return nil
}

func (r *FileRepository) Ping() error {
	return nil
}

func (r *FileRepository) AddBatch(data []map[string]string) ([]byte, error) {
	var res []urlBatchData

	file, err := os.OpenFile(r.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	m := make(map[string]string)
	for _, v := range data {
		code := v["correlation_id"]
		original := v["original_url"]
		short := services.MakeShortURL(original)

		m[short] = original

		res = append(res, urlBatchData{
			ID:    code,
			Short: short,
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

	responseJSON, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	return responseJSON, nil
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
