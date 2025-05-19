package memoryrepository

import (
	"context"

	"github.com/kirsh-nat/shortener.git/internal/services"
)

// TODO: add user
func (r *MemoryRepository) AddBatch(context context.Context, host, userID string, data []services.BatchItem) ([]services.URLData, error) {
	var res []services.URLData

	for _, v := range data {
		short := services.MakeShortURL(v.Original)
		err := r.Add(context, services.MakeFullShortURL(short, host), v.Original, userID)
		if err != nil {
			return nil, err
		}

		res = append(res, services.URLData{
			ID:    v.ID,
			Short: services.MakeFullShortURL(short, host),
		})
	}

	return res, nil
}
