package dbrepository

import (
	"context"
	"database/sql"

	"github.com/kirsh-nat/shortener.git/internal/domain"
)

func (r *DBRepository) Get(context context.Context, short string) (string, error) {
	row := r.db.QueryRowContext(context,
		"SELECT original_url from links where short_url = $1", short)
	var long sql.NullString

	err := row.Scan(&long)
	if err != nil {
		return "", err
	}
	if long.Valid {
		return long.String, nil
	}

	return "", domain.ErrorURLNotFound
}
