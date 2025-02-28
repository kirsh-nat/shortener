package main

import (
	"net/http"

	"github.com/kirsh-nat/shortener.git/internal/config"
)

var (
	URLList = make(map[string]string)
	conf    = new(config.Config)
)

func init() {
	config.SetConfig(conf)
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	mux := routes()
	return http.ListenAndServe(conf.Addr, mux)
}
