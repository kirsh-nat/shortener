package main

import (
	"net/http"

	"github.com/kirsh-nat/shortener.git/internal/app"
)

func main() {
	app.SetAppConfig()
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	mux := app.Routes()
	return http.ListenAndServe(app.AppSettings.Addr, mux)
}
