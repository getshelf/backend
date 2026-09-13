package values

import (
	"regexp"

	domainErrors "github.com/getshelf/backend/internal/domain/user/errors"
)

var usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9_]{3,32}$`)

type Username string

func NewUsername(raw string) (Username, error) {
	if !usernamePattern.MatchString(raw) {
		return "", domainErrors.ErrInvalidUsername
	}
	return Username(raw), nil
}

func (username Username) String() string {
	return string(username)
}
