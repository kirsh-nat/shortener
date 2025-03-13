package app

import (
	"github.com/kirsh-nat/shortener.git/internal/config"
)

var (
	listURL     = make(map[string]string)
	AppSettings = new(config.Config)
)

func SetAppConfig() {
	AppSettings = new(config.Config)
	config.ParseFlags(AppSettings)
	config.ValidateConfig(AppSettings)
}
