package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"regexp"
)

const (
	ErrURLNotFound = "URL not found"
	ErrURLExist    = "URL already exists"
	ErrURLFormat   = "string is not a valid URL"
)

func makeShortURL(url string) string {
	hash := sha256.New()
	hash.Write([]byte(url))
	hashed := hash.Sum(nil)

	return hex.EncodeToString(hashed)[:6]
}

func AddURL(longURL string, listURL *map[string]string) (string, error) {
	err := validateLongURL(longURL)
	if err != nil {
		return "", err
	}

	shortURL := makeShortURL(longURL)

	if _, exists := (*listURL)[shortURL]; exists {
		return "", errors.New(ErrURLExist)
	}

	(*listURL)[shortURL] = longURL

	return shortURL, nil
}

func GetURL(shortURL string, listURL map[string]string) (string, error) {
	redirectURL, exists := listURL[shortURL]
	if !exists {
		return "", errors.New(ErrURLNotFound)
	}

	return redirectURL, nil
}

func validateLongURL(longURL string) error {
	var validURL = regexp.MustCompile(`^(http|https):\/\/[a-zA-Z0-9\-]+\.[a-zA-Z]{2,6}(\/\S*)?$`)
	if !validURL.MatchString(longURL) {
		return errors.New(ErrURLFormat)
	}

	return nil
}
