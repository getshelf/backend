package account

import (
	"context"
	"time"
)

type CreateAccountParams struct {
	ID           string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Store interface {
	Create(
		ctx context.Context,
		params CreateAccountParams,
	) error
	FindByEmail(
		ctx context.Context,
		email string,
	) (Account, error) // replace with Account
}
