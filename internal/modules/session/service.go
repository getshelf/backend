package session

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Service struct {
	store Store
	ttl   time.Duration
	now   func() time.Time
}

func NewService(store Store, ttl time.Duration) *Service {
	return &Service{
		store: store,
		ttl:   ttl,
		now:   time.Now,
	}
}

func (service *Service) Create(
	ctx context.Context,
	accountID string,
) (CreatedSession, error) {
	token, tokenHash, err := newSessionToken()
	if err != nil {
		return CreatedSession{}, fmt.Errorf(
			"generate session token: %w",
			err,
		)
	}

	now := service.now().UTC()
	expiresAt := now.Add(service.ttl)

	err = service.store.Create(
		ctx,
		CreateSessionParams{
			TokenHash: tokenHash,
			AccountID: accountID,
			ExpiresAt: expiresAt,
			CreatedAt: now,
		},
	)
	if err != nil {
		return CreatedSession{}, fmt.Errorf(
			"create session: %w",
			err,
		)
	}

	return CreatedSession{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func (service *Service) Authenticate(
	ctx context.Context,
	token string,
) (string, error) {
	if token == "" {
		return "", ErrInvalidSession
	}

	found, err := service.store.FindByTokenHash(
		ctx,
		hashSessionToken(token),
	)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return "", ErrInvalidSession
		}

		return "", fmt.Errorf("find session: %w", err)
	}

	if !found.ExpiresAt.After(service.now()) {
		return "", ErrInvalidSession
	}

	return found.AccountID, nil
}


//logout 
func (service *Service) Logout(ctx context.Context, token string) error {
	if token == "" {
		return nil
	}
	
	if err := service.store.DeleteByTokenHash(ctx, hashSessionToken(token)); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	return nil;
}