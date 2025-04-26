package dbrepository

import (
	"context"

	"github.com/kirsh-nat/shortener.git/internal/services"
)

func (r *DBRepository) AddBatch(context context.Context, host string, data []services.BatchItem) ([]services.URLData, error) {
	var res []services.URLData

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(context,
		"INSERT INTO links (short_url, original_url) VALUES($1, $2)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	for _, v := range data {
		short := services.MakeShortURL(v.Original)

		_, err := stmt.ExecContext(context, short, v.Original)
		if err != nil {
			return nil, err
		}

		res = append(res, services.URLData{
			ID:    v.ID,
			Short: services.MakeFullShortURL(short, host),
		})
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return res, nil
}
