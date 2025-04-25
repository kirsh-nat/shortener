package main

import (
	"net/http"

	"github.com/kirsh-nat/shortener.git/internal/app"
	"github.com/kirsh-nat/shortener.git/internal/handlers"
	"github.com/kirsh-nat/shortener.git/internal/repositories"
	"github.com/kirsh-nat/shortener.git/internal/services"
)

func main() {
	app.SetAppConfig()

	app.ParseFlags(app.AppSettings)
	app.ValidateConfig(app.AppSettings)

	var repo services.URLRepository

	if app.AppSettings.SetDBConnection != "" {
		app.DB = app.DBConnect(app.AppSettings.SetDBConnection)
		repo = repositories.NewDBRepository(app.DB)

	} else if app.AppSettings.FilePath != "" {
		repo = repositories.NewFileRepository(app.AppSettings.FilePath)

	} else {
		repo = repositories.NewMemoryRepository()

	}
	service := services.NewURLService(repo)
	handler := handlers.NewURLHandler(service)

	if err := run(handler); err != nil {
		app.Sugar.Fatalw(err.Error(), "event", "start server")
	}
	if app.DB != nil {
		defer app.DB.Close()

	}
}

func run(handler *handlers.URLHandler) error {
	mux := handlers.Routes(handler)
	return http.ListenAndServe(app.AppSettings.Addr, mux)
}
