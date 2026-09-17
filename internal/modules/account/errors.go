package account

import "errors"

var (
	ErrInvalidEmail          = errors.New("invalid email")
	ErrPasswordTooWeak       = errors.New("password must contain at least eight characters")
	ErrPasswordContainsEmail = errors.New("password must not contain the email")
	ErrEmailTaken            = errors.New("email is already taken")
)
