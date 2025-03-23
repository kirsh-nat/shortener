package main

import (
	"net/http"

	"github.com/kirsh-nat/shortener.git/internal/app"
	"github.com/kirsh-nat/shortener.git/internal/config"
)

func main() {
	app.SetAppConfig()

	config.ParseFlags(app.AppSettings)
	config.ValidateConfig(app.AppSettings)

	app.Store = app.NewURLStore(app.AppSettings.FilePath)

	if err := run(); err != nil {
		app.Sugar.Fatalw(err.Error(), "event", "start server")
		//panic(err)
	}
}

func run() error {
	mux := app.Routes()
	return http.ListenAndServe(app.AppSettings.Addr, mux)
}
