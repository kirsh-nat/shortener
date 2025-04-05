package app

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kirsh-nat/shortener.git/internal/config"
	"go.uber.org/zap"
)

var (
	AppSettings = new(config.Config)
	Store       *URLStore
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

func setLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		Sugar.Fatalw(err.Error(), "event", err)

	}
	defer logger.Sync()

	Sugar = *logger.Sugar()
}

func SetDBConnection(ps string) *sql.DB {
	DB, err := sql.Open("pgx", ps)
	if err != nil {
		Sugar.Fatalw(err.Error(), "event", err)
		return nil
	}

	return DB
	//defer db.Close()
	// ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	// defer cancel()
	// if err = db.PingContext(ctx); err != nil {
	// 	panic(err)
	// }
}
