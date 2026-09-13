package entities

import (
	"context"
	"time"

	domainErrors "github.com/getshelf/backend/internal/domain/user/errors"
	"github.com/getshelf/backend/internal/domain/user/ports"
	"github.com/getshelf/backend/internal/domain/user/values"
)

type User struct {
	id           values.UserID
	email        values.Email
	passwordHash values.PasswordHash
	createdAt    time.Time
}

func NewUser(ctx context.Context, id values.UserID, email values.Email, hash values.PasswordHash, checker ports.UniquenessChecker, now time.Time) (User, error) {
	validatedEmail, err := values.NewEmail(email.String())
	if err != nil {
		return User{}, err
	}
	validatedHash, err := values.NewPasswordHash(hash.String())
	if err != nil {
		return User{}, err
	}

	emailTaken, err := checker.IsEmailTaken(ctx, validatedEmail)
	if err != nil {
		return User{}, err
	}
	if emailTaken {
		return User{}, domainErrors.ErrEmailTaken
	}

	return User{id: id, email: validatedEmail, passwordHash: validatedHash, createdAt: now}, nil
}

func (user User) ID() values.UserID { return user.id }

func (user User) Email() values.Email { return user.email }

func (user User) PasswordHash() values.PasswordHash { return user.passwordHash }

func (user User) CreatedAt() time.Time { return user.createdAt }
