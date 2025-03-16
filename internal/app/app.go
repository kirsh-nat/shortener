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
	AppSettings = new(config.Config)
	Store = NewURLStore()
	config.ParseFlags(AppSettings)
	config.ValidateConfig(AppSettings)

	Sugar.Infow(
		"Starting server",
		"addr", AppSettings.Addr,
	)

}
