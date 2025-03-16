package main

import (
	"net/http"

	"github.com/kirsh-nat/shortener.git/internal/app"
	"go.uber.org/zap"
)

// var sugar zap.SugaredLogger

func main() {
	// создаём предустановленный регистратор zap
	logger, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic(err)
	}
	defer logger.Sync()

	// делаем регистратор SugaredLogger
	// TODO: возмлонжно вообще вынести в конфигу аппы????
	app.Sugar = *logger.Sugar()

	app.SetAppConfig()

	if err := run(); err != nil {
		app.Sugar.Fatalw(err.Error(), "event", "start server")
		//panic(err)
	}
}

func run() error {
	mux := app.Routes()
	return http.ListenAndServe(app.AppSettings.Addr, mux)
}
