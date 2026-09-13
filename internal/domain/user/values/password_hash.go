package values

import (
	"strings"

	domainErrors "github.com/getshelf/backend/internal/domain/user/errors"
)

type PasswordHash string

func NewPasswordHash(raw string) (PasswordHash, error) {
	if strings.TrimSpace(raw) == "" {
		return "", domainErrors.ErrInvalidPasswordHash
	}
	return PasswordHash(raw), nil
}

func (hash PasswordHash) String() string {
	return string(hash)
}
