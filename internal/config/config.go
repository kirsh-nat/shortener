package config

import (
	"flag"
)

type Config struct {
	Addr string
	Resp string
}

func SetConfig(c *Config) {
	flag.StringVar(&c.Addr,
		"a", "localhost:8080",
		"Адрес запуска HTTP-сервера",
	)
	flag.StringVar(&c.Resp,
		"b", "localhost:8080",
		"Базовый адрес результирующего сокращённого URL ",
	)

}
