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

//TODO: метод на объявление логгера и записи логов в файлик!!!!

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
