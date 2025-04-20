package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kirsh-nat/shortener.git/internal/domain"
	"github.com/kirsh-nat/shortener.git/internal/migrations"
	"github.com/kirsh-nat/shortener.git/internal/models"
	"github.com/kirsh-nat/shortener.git/internal/services"
)

type DBRepository struct {
	db *sql.DB
}

func NewDBRepository(db *sql.DB) models.URLRepository {
	migrations.CreateLinkTable(db)
	return &DBRepository{db: db}
}

func (r *DBRepository) Add(shortURL, originalURL string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
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

func (r *DBRepository) Get(short string) (string, error) {
	row := r.db.QueryRowContext(context.Background(),
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

func (r *DBRepository) AddBatch(host string, data []map[string]string) ([]byte, error) {
	type urlData struct {
		ID    string `json:"correlation_id"`
		Short string `json:"short_url"`
	}

	var res []urlData

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stmt, err := tx.PrepareContext(ctx,
		"INSERT INTO links (short_url, original_url) VALUES($1, $2)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	for _, v := range data {
		code := v["correlation_id"]
		original := v["original_url"]
		short := services.MakeShortURL(original)

		_, err := stmt.ExecContext(ctx, short, original)
		if err != nil {
			return nil, err
		}

		res = append(res, urlData{
			ID:    code,
			Short: "http://" + host + "/" + short,
		})
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	responseJSON, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	return responseJSON, nil
}
