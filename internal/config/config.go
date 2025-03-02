package config

import (
	"flag"
	"os"
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

func ValidateConfig(c *Config) {
	if envAddr := os.Getenv("SERVER_ADDRESS"); envAddr != "" {
		c.Addr = envAddr
	}
	if envResp := os.Getenv("WEB_ADDRESS"); envResp != "" {
		c.Resp = envResp
	}
	if c.Resp != c.Addr {
		c.Resp = c.Addr
	}
}
