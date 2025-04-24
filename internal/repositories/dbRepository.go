package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kirsh-nat/shortener.git/internal/app"
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

func (r *DBRepository) Add(ctx context.Context, shortURL, originalURL, userID string) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO links (short_url, original_url, user_id, deleted) VALUES ($1, $2, $3, $4)", shortURL, originalURL, userID, false)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return domain.NewDublicateError("DB dublicate error", err)
		}
		return err
	}
	r.AddUserURL(userID, shortURL)

	return nil

}

func (r *DBRepository) Get(short string) (string, error) {
	row := r.db.QueryRowContext(context.Background(),
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
			return long.String, fmt.Errorf("url deleted")
		}
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

func (r *DBRepository) DeleteBatch(shortURLs []string, userID string) {
	const batchSize = 5
	var wg sync.WaitGroup
	errChan := make(chan error, len(shortURLs))
	urlChan := make(chan string)

	go func() {
		for err := range errChan {
			if err != nil {
				app.Sugar.Info("Error:", err)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for url := range urlChan {
			if err := r.updateDeletedStatus(url, userID); err != nil {
				errChan <- err
			}
		}
	}()

	for i := 0; i < len(shortURLs); i += batchSize {
		end := i + batchSize
		if end > len(shortURLs) {
			end = len(shortURLs)
		}

		batch := shortURLs[i:end]

		for _, url := range batch {
			urlChan <- url
		}

	}

	close(urlChan)
	wg.Wait()
	close(errChan)
}

func (r *DBRepository) updateDeletedStatus(shortURL, userID string) error {
	ctx := context.Background()
	result, err := r.db.ExecContext(ctx,
		"UPDATE links SET deleted = $1 WHERE short_url = $2 AND user_id = $3", true, shortURL, userID)

	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no rows updated for URL: %s", shortURL)
	}

	return nil
}

func (r *DBRepository) AddUserURL(userID, short string) {
}

func (r *DBRepository) GetUserURLs(userID string) ([]string, error) {
	rows, err := r.db.QueryContext(context.Background(),
		"SELECT short_url FROM links WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}
