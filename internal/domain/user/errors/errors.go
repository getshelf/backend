package errors

import "errors"

var (
	ErrInvalidUsername     = errors.New("invalid username")
	ErrInvalidEmail        = errors.New("invalid email")
	ErrInvalidPasswordHash = errors.New("invalid password hash")
	ErrPasswordTooWeak     = errors.New("password is too weak")
	ErrUsernameTaken       = errors.New("username is already taken")
	ErrEmailTaken          = errors.New("email is already taken")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrUserNotFound        = errors.New("user not found")
)
