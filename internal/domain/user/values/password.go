package values

import (
	"strings"
	"unicode"

	domainErrors "github.com/getshelf/backend/internal/domain/user/errors"
)

type PlainPassword string

func NewPlainPassword(password string, email Email) (PlainPassword, error) {
	if len([]rune(password)) < 8 || !hasLetter(password) || !hasDigit(password) {
		return "", domainErrors.ErrPasswordTooWeak
	}

	passwordLower := strings.ToLower(password)
	if strings.Contains(passwordLower, strings.ToLower(email.String())) {
		return "", domainErrors.ErrPasswordTooWeak
	}
	return PlainPassword(password), nil
}

func (password PlainPassword) String() string {
	return string(password)
}

func hasLetter(value string) bool {
	for _, character := range value {
		if unicode.IsLetter(character) {
			return true
		}
	}
	return false
}

func hasDigit(value string) bool {
	for _, character := range value {
		if unicode.IsDigit(character) {
			return true
		}
	}
	return false
}
