package app

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

var (
	AppSettings *Config
	Sugar       zap.SugaredLogger
	DB          *sql.DB
)

func SetAppConfig() {
	setLogger()
	AppSettings = new(Config)
	Sugar.Infow(
		"Starting server",
		"addr", AppSettings.Addr,
	)

}
