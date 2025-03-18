package app

import (
	"github.com/kirsh-nat/shortener.git/internal/config"
	"go.uber.org/zap"
)

var (
	AppSettings = new(config.Config)
	Store       *URLStore
	Sugar       zap.SugaredLogger
)

func SetAppConfig() {
	setLogger()
	AppSettings = new(config.Config)
	Store = NewURLStore()
	config.ParseFlags(AppSettings)
	config.ValidateConfig(AppSettings)

	Sugar.Infow(
		"Starting server",
		"addr", AppSettings.Addr,
	)

}

func setLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	Sugar = *logger.Sugar()
}
