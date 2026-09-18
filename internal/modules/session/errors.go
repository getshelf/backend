package session

import "errors"

var (
	ErrInvalidSession          = errors.New("invalid session")
	ErrNotFound                = errors.New("not found")
)
