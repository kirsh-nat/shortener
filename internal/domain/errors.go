package domain

import (
	"errors"
	"fmt"
	"strings"
)

var (
	errorExistURL    = errors.New("URL already exists")
	ErrorURLNotFound = errors.New("URL not found")

//	errorURL    = errors.New("URL not found")

// errorUnknownStorage = errors.New("unknown storage type")
)

type DublicateError struct {
	level string
	Err   error
}

func (le *DublicateError) Error() string {
	return fmt.Sprintf("[%s] %v", le.level, le.Err)
}

func NewDublicateError(label string, err error) error {
	return &DublicateError{
		level: strings.ToUpper(label),
		Err:   err,
	}
}
