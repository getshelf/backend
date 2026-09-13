package errors

import "errors"

var (
	ErrInvalidEmail        = errors.New("invalid email")
	ErrInvalidPasswordHash = errors.New("invalid password hash")
	ErrPasswordTooWeak     = errors.New("password is too weak")
	ErrEmailTaken          = errors.New("email is already taken")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrUserNotFound        = errors.New("user not found")
)
