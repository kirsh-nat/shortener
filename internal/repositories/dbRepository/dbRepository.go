package dbrepository

import (
	"database/sql"

	"github.com/kirsh-nat/shortener.git/cmd/shortener/migrations"
	"github.com/kirsh-nat/shortener.git/internal/services"
)

type DBRepository struct {
	db *sql.DB
}

func NewDBRepository(db *sql.DB) services.URLRepository {
	migrations.CreateLinkTable(db)
	return &DBRepository{db: db}
}

func (r *DBRepository) Ping() error {
	if err := r.db.Ping(); err != nil {
		return err
	}

	return nil
}
