package app

import (
	"github.com/kirsh-nat/shortener.git/internal/config"
)

var (
	AppSettings = new(config.Config)
	Store       *URLStore
)

func SetAppConfig() {
	AppSettings = new(config.Config)
	Store = NewURLStore()
	config.ParseFlags(AppSettings)
	config.ValidateConfig(AppSettings)
}
