package ports

import (
	"context"
	"time"

	"github.com/getshelf/backend/internal/domain/user/events"
	"github.com/getshelf/backend/internal/domain/user/values"
)

type UniquenessChecker interface {
	IsEmailTaken(context.Context, values.Email) (bool, error)
}

type User interface {
	ID() values.UserID
	Email() values.Email
	PasswordHash() values.PasswordHash
	CreatedAt() time.Time
}

type UserRepository interface {
	Save(context.Context, User) error
}

type AuthUser struct {
	ID           values.UserID
	Email        values.Email
	PasswordHash values.PasswordHash
	CreatedAt    time.Time
}

type AuthRepository interface {
	FindByLogin(context.Context, string) (AuthUser, error)
	FindByID(context.Context, values.UserID) (AuthUser, error)
	UpdateProfile(context.Context, values.UserID, values.Email) error
}

type PasswordHasher interface {
	Hash(context.Context, values.PlainPassword) (values.PasswordHash, error)
}

type PasswordVerifier interface {
	Compare(values.PasswordHash, string) error
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type TokenClaims struct {
	UserID    values.UserID
	TokenID   string
	ExpiresAt time.Time
	Refresh   bool
}

type TokenService interface {
	Issue(AuthUser) (TokenPair, error)
	Parse(string, bool) (TokenClaims, error)
}

type TokenBlacklist interface {
	Revoke(context.Context, string, time.Duration) error
	IsRevoked(context.Context, string) (bool, error)
}

type IDGenerator interface {
	NewUserID(context.Context) (values.UserID, error)
}

type Clock interface {
	Now() time.Time
}

type TxManager interface {
	Do(context.Context, func(context.Context) error) error
}

type EventPublisher interface {
	Publish(context.Context, events.UserRegistered) error
}
