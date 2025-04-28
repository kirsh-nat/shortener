package dbrepository

import (
	"context"
	"database/sql"

	"github.com/kirsh-nat/shortener.git/internal/domain"
)

func (r *DBRepository) Get(context context.Context, short string) (string, error) {
	row := r.db.QueryRowContext(context,
		"SELECT original_url, deleted FROM links WHERE short_url = $1", short)

	var long sql.NullString
	var deleted bool

	err := row.Scan(&long, &deleted)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", domain.ErrorURLNotFound
		}
		return "", err
	}

	if long.Valid {
		if deleted {
			return long.String, domain.NewDeletedError("URL deleted", nil)
		}
		return long.String, nil
	}

	return "", domain.ErrorURLNotFound
}
