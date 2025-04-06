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
	app.DB = app.SetDBConnection(app.AppSettings.SetDBConnection)

	app.Store = app.NewURLStore(app.AppSettings)

	if err := run(); err != nil {
		app.Sugar.Fatalw(err.Error(), "event", "start server")
	}
	if app.Store.DBConnection != nil {
		defer app.Store.DBConnection.Close()

	}
}

func run() error {
	mux := app.Routes()
	return http.ListenAndServe(app.AppSettings.Addr, mux)
}
