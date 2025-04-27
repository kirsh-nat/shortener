package app

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"github.com/kirsh-nat/shortener.git/internal/config"
)

var (
	AppSettings *config.Config
	Sugar       zap.SugaredLogger
	DB          *sql.DB
)

func SetAppConfig() {
	setLogger()
	AppSettings = new(config.Config)
	Sugar.Infow(
		"Starting server",
		"addr", AppSettings.Addr,
	)

}
