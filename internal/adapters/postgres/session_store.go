package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/getshelf/backend/internal/adapters/postgres/sqlcgen"
	"github.com/getshelf/backend/internal/modules/session"
)

type SessionStore struct {
	queries *sqlcgen.Queries
}

func NewSessionStore(db *sql.DB) *SessionStore {
	return &SessionStore{
		queries: sqlcgen.New(db),
	}
}

func (store *SessionStore) Create (
	ctx context.Context,
	params session.CreateSessionParams,
) error {
	err := store.queries.CreateSession(
		ctx,
		sqlcgen.CreateSessionParams{
			TokenHash: params.TokenHash[:],
			AccountID: params.AccountID,
			ExpiresAt: params.ExpiresAt,
			CreatedAt: params.CreatedAt,
		},
	)
	return err
}

func (store *SessionStore) FindByTokenHash (
	ctx context.Context,
	tokenHash [32]byte,
) (session.Session, error) {
	s, err := store.queries.FindSessionByTokenHash(
		ctx,
		tokenHash[:],
	)

	if err != nil {
		return session.Session{}, err
	}

	return session.Session{
		TokenHash: s.TokenHash,
		AccountID: s.AccountID,
		ExpiresAt: s.ExpiresAt,
		CreatedAt: s.CreatedAt,
	}, nil
}

func (store *SessionStore) DeleteByTokenHash(
	ctx context.Context,
	tokenHash [32]byte,
) error {
	return store.queries.DeleteSessionByTokenHash(
		ctx,
		tokenHash[:],
	)
}

func (store *SessionStore) DeleteExpired(
	ctx context.Context,
	before time.Time,
) error {
	return store.queries.DeleteExpiredSessions(
		ctx,
		before,
	)
}
