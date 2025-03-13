package config

import (
	"flag"
	"fmt"
	"os"
)

const (
	defAddr    string = "localhost:8080" // Дефолтный адрес запуска HTTP-сервера
	defResp    string = "localhost:8080" // Дефолтный базовый адрес результирующего сокращённого URL
	srvAddrVar string = "SERVER_ADDRESS" // Переменная окружения для адреса запуска HTTP-сервера
	webAddrVar string = "WEB_ADDRESS"    // Переменная окружения для базового адреса результирующего сокращённого URL
)

type Config struct {
	Addr string
	Resp string
}

func ValidateConfig(c *Config) {
	fmt.Print(c)
	// Пока у нас только один сервер, запросы на короткий урл должны ссылаться также на него
	if c.Resp != c.Addr {
		c.Resp = c.Addr
	}
}

func ParseFlags(c *Config) {
	flag.StringVar(&c.Addr,
		"a", defAddr,
		"Адрес запуска HTTP-сервера",
	)
	flag.StringVar(&c.Resp,
		"b", defResp,
		"Базовый адрес результирующего сокращённого URL ",
	)
	flag.Parse()

	//Если заданы переменные окружения, меняем настройки в соответвии с ними
	if envAddr := os.Getenv(srvAddrVar); envAddr != "" {
		c.Addr = envAddr
	}
	if envResp := os.Getenv(webAddrVar); envResp != "" {
		c.Resp = envResp
	}
}
