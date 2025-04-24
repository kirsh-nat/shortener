package app

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

var (
	AppSettings *Config
	//	Store       *URLStore
	Sugar zap.SugaredLogger
	DB    *sql.DB
)

func SetAppConfig() {
	setLogger()
	AppSettings = new(Config)
	Sugar.Infow(
		"Starting server",
		"addr", AppSettings.Addr,
	)

}

// TODO разные файлы
func setLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		Sugar.Fatalw(err.Error(), "event", err)

	}
	defer logger.Sync()

	Sugar = *logger.Sugar()
}

func DBConnect(ps string) *sql.DB {
	DB, err := sql.Open("pgx", ps)
	if err != nil {
		Sugar.Fatalw(err.Error(), "event", err)
		return nil
	}

	if err := DB.Ping(); err != nil {
		Sugar.Fatalw(err.Error(), "event", err)
		return nil
	}

	return DB
}
