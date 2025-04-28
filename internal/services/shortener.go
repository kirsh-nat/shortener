package services

import (
	"crypto/sha256"
	"encoding/hex"
)

const (
	lenURL = 6 // length of the short URL
)

func MakeShortURL(url string) string {
	hash := sha256.New()
	hash.Write([]byte(url))
	hashed := hash.Sum(nil)

	return hex.EncodeToString(hashed)[:lenURL]
}

func MakeFullShortURL(codeURL, host string) string {
	return "http://" + host + "/" + codeURL
}
