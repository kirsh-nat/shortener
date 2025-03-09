package internal

import (
	"crypto/sha256"
	"encoding/hex"
)

func NewShortURL(url string) string {
	hash := sha256.New()
	hash.Write([]byte(url))
	hashed := hash.Sum(nil)

	return hex.EncodeToString(hashed)[:6]
}
