package db

import (
	"database/sql"

	"go.uber.org/zap"
)

func DBConnect(ps string, logger zap.SugaredLogger) *sql.DB {
	DB, err := sql.Open("pgx", ps)
	if err != nil {
		logger.Fatalw(err.Error(), "event", err)
		return nil
	}

	return DB
}
