package main

import (
	"net/http"
)

var UrlList = make(map[string]string)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {

	mux := routes()
	return http.ListenAndServe(":8080", mux)
}
