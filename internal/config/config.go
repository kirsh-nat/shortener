package config

import (
	"flag"
	"os"
)

const (
	defAddr         string = "localhost:8080"    // Дефолтный адрес запуска HTTP-сервера
	defResp         string = "localhost:8080"    // Дефолтный базовый адрес результирующего сокращённого URL
	defPath         string = "/tmp/urls.txt"     // Дефолтный адрес файла с ссылками
	defPathVar      string = "FILE_STORAGE_PATH" // Переменная окружения для файла с ссылками
	srvAddrVar      string = "SERVER_ADDRESS"    // Переменная окружения для адреса запуска HTTP-сервера
	webAddrVar      string = "WEB_ADDRESS"       // Переменная окружения для базового адреса результирующего сокращённого URL
	SetDBConnection string = "DATABASE_DSN"      // Переменная окружения для базового адреса результирующего сокращённого URL
)

type Config struct {
	Addr            string
	Resp            string
	FilePath        string
	SetDBConnection string
}

func ValidateConfig(c *Config) {
	if c.Addr == "" {
		c.Addr = "localhost:8080"
	}
	// Пока у нас только один сервер, запросы на короткий урл должны ссылаться также на него
	if c.Resp == "" || c.Resp != c.Addr {
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

	flag.StringVar(&c.FilePath,
		"f", "",
		"Путь к файлу с ссылками",
	)

	flag.StringVar(&c.SetDBConnection,
		"d", "",
		"Строка подключения к базе данных",
	)

	flag.Parse()

	if envAddr := os.Getenv(srvAddrVar); envAddr != "" {
		c.Addr = envAddr
	}
	if envResp := os.Getenv(webAddrVar); envResp != "" {
		c.Resp = envResp
	}
	if envPath := os.Getenv(defPathVar); envPath != "" {
		c.FilePath = envPath
	}
	if envPath := os.Getenv(SetDBConnection); envPath != "" {
		c.SetDBConnection = envPath
	}
}
