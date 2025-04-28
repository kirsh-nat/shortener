package dbrepository

import (
	"context"

	"github.com/kirsh-nat/shortener.git/internal/services"
)

func (r *DBRepository) GetUserURLs(ctx context.Context, userID string) ([]services.UserURLData, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT short_url as Short, original_url as Original FROM links WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userUrls []services.UserURLData
	for rows.Next() {
		var urlData services.UserURLData
		if err := rows.Scan(&urlData.Short, &urlData.Original); err != nil {
			return nil, err
		}
		userUrls = append(userUrls, urlData)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userUrls, nil
}
