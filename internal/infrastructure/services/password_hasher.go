package services

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/getshelf/backend/internal/domain/user/values"
)

type BcryptPasswordHasher struct {
	cost int
}

func NewBcryptPasswordHasher(cost int) BcryptPasswordHasher {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	return BcryptPasswordHasher{cost: cost}
}

func (hasher BcryptPasswordHasher) Hash(_ context.Context, password values.PlainPassword) (values.PasswordHash, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password.String()), hasher.cost)
	if err != nil {
		return "", err
	}
	return values.NewPasswordHash(string(hash))
}

func (BcryptPasswordHasher) Compare(hash values.PasswordHash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash.String()), []byte(password))
}
