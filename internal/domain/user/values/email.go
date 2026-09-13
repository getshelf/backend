package values

import (
	"regexp"
	"strings"

	domainErrors "github.com/getshelf/backend/internal/domain/user/errors"
)

var emailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type Email string

func NewEmail(email string) (Email, error) {
	normalized := strings.ToLower(strings.TrimSpace(email))
	if !emailPattern.MatchString(normalized) {
		return "", domainErrors.ErrInvalidEmail
	}
	return Email(normalized), nil
}

func (email Email) String() string {
	return string(email)
}
