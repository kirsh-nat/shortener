package domain

import (
	"errors"
	"fmt"
	"strings"
)

type DublicateError struct {
	level string
	Err   error
}
type DeletedError struct {
	level string
	Err   error
}

var (
	ErrorURLNotFound = errors.New("URL not found")
	ErrorURLDeleted  = errors.New("URL nwas deleted")
)

func (le *DublicateError) Error() string {
	return fmt.Sprintf("[%s] %v", le.level, le.Err)
}

func NewDublicateError(label string, err error) error {
	return &DublicateError{
		level: strings.ToUpper(label),
		Err:   err,
	}
}

func (le *DeletedError) Error() string {
	return fmt.Sprintf("[%s] %v", le.level, le.Err)
}

func NewDeletedError(label string, err error) error {
	return &DeletedError{
		level: strings.ToUpper(label),
		Err:   err,
	}
}
