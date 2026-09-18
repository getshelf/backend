package session

import (
	"context"
	"time"
)

type CreateSessionParams struct {
	TokenHash [32]byte
	AccountID string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type Store interface {
	Create(
		ctx context.Context,
		params CreateSessionParams,
	) error

	FindByTokenHash(
		ctx context.Context,
		tokenHash [32]byte,
	) (Session, error)

	DeleteByTokenHash(
		ctx context.Context,
		tokenHash [32]byte,
	) error

	DeleteExpired(
		ctx context.Context,
		before time.Time,
	) error
}
