package main

import (
	"net/http"

	"github.com/kirsh-nat/shortener.git/internal/app"
	"github.com/kirsh-nat/shortener.git/internal/config"
	"github.com/kirsh-nat/shortener.git/internal/db"
	"github.com/kirsh-nat/shortener.git/internal/handlers"
	dbrepository "github.com/kirsh-nat/shortener.git/internal/repositories/dbRepository"
	filerepository "github.com/kirsh-nat/shortener.git/internal/repositories/fileRepository"
	memoryrepository "github.com/kirsh-nat/shortener.git/internal/repositories/memoryRepository"
	"github.com/kirsh-nat/shortener.git/internal/services"
)

func main() {
	app.SetAppConfig()

	config.ParseFlags(app.AppSettings)
	config.ValidateConfig(app.AppSettings)

	var repo services.URLRepository

	if app.AppSettings.SetDBConnection != "" {
		app.DB = db.DBConnect(app.AppSettings.SetDBConnection, app.Sugar)
		repo = dbrepository.NewDBRepository(app.DB)

	} else if app.AppSettings.FilePath != "" {
		repo = filerepository.NewFileRepository(app.AppSettings.FilePath)

	} else {
		repo = memoryrepository.NewMemoryRepository()

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
