package repositories

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kirsh-nat/shortener.git/cmd/shortener/migrations"
	"github.com/kirsh-nat/shortener.git/internal/domain"
	"github.com/kirsh-nat/shortener.git/internal/services"
)

type DBRepository struct {
	db *sql.DB
}

func NewDBRepository(db *sql.DB) services.URLRepository {
	migrations.CreateLinkTable(db)
	return &DBRepository{db: db}
}

func (r *DBRepository) Add(ctx context.Context, shortURL, originalURL string) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO links (short_url, original_url) VALUES ($1, $2)", shortURL, originalURL)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return domain.NewDublicateError("DB dublicate error", err)
		}
		return err
	}

	return nil

}

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

func (r *DBRepository) Ping() error {
	if err := r.db.Ping(); err != nil {
		return err
	}

	return nil
}

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
			Short: "http://" + host + "/" + short,
		})
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return res, nil
}
