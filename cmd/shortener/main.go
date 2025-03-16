package main

import (
	"net/http"

	"github.com/kirsh-nat/shortener.git/internal/app"
)

// var sugar zap.SugaredLogger

func main() {
	// создаём предустановленный регистратор zap

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
