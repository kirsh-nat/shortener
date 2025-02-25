package internal

import (
	"math/rand"
	"time"
)

const (
	availableChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	maxLen         = 6
)

func NewShortUrl(url string) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune(availableChars)

	short := make([]rune, maxLen)
	for i := range short {
		short[i] = chars[rnd.Intn(len(chars))]
	}

	return string(short)
}
