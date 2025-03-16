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
		// вызываем панику, если ошибка
		panic(err)
	}
	defer logger.Sync()

	// делаем регистратор SugaredLogger
	// TODO: возмлонжно вообще вынести в конфигу аппы????
	Sugar = *logger.Sugar()
}
