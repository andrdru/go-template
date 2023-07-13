package entities

import (
	"errors"
)

var (
	ErrNotFound = errors.New("not found")

	ErrNotAllowed = errors.New("not allowed")

	ErrAlreadyExists = errors.New("already exists")

	errInternal = errors.New("internal error")
)

func Err(err error) func() string {
	return func() string {
		if err == nil {
			return ""
		}

		switch true {
		case errors.Is(err, ErrNotFound):
			err = ErrNotFound
		case errors.Is(err, ErrNotAllowed):
			err = ErrNotAllowed
		case errors.Is(err, ErrAlreadyExists):
			err = ErrAlreadyExists
		default:
			err = errInternal
		}

		return err.Error()
	}
}
