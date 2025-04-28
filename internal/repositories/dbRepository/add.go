package dbrepository

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kirsh-nat/shortener.git/internal/domain"
)

func (r *DBRepository) Add(ctx context.Context, shortURL, originalURL, userID string) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO links (short_url, original_url, user_id, deleted) VALUES ($1, $2, $3, $4)", shortURL, originalURL, userID, false)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return domain.NewDublicateError("DB dublicate error", err)
		}
		return err
	}

	return nil
}
