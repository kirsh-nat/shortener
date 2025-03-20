package main

import (
	"net/http"

	"github.com/kirsh-nat/shortener.git/internal/app"
	"github.com/kirsh-nat/shortener.git/internal/config"
)

// var sugar zap.SugaredLogger

func main() {
	// создаём предустановленный регистратор zap

	app.SetAppConfig()
	config.ParseFlags(app.AppSettings)
	config.ValidateConfig(app.AppSettings)

	if err := run(); err != nil {
		app.Sugar.Fatalw(err.Error(), "event", "start server")
		//panic(err)
	}
}

func run() error {
	mux := app.Routes()
	return http.ListenAndServe(app.AppSettings.Addr, mux)
}
