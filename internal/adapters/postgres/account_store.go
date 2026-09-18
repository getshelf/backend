package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/getshelf/backend/internal/adapters/postgres/sqlcgen"
	"github.com/getshelf/backend/internal/modules/account"
	"github.com/lib/pq"
)

type AccountStore struct {
	queries *sqlcgen.Queries
}

func NewAccountStore(db *sql.DB) *AccountStore {
	return &AccountStore{
		queries: sqlcgen.New(db),
	}
}

func (store *AccountStore) Create(
	ctx context.Context,
	params account.CreateAccountParams,
) error {
	err := store.queries.CreateAccount(
		ctx,
		sqlcgen.CreateAccountParams{
			ID:           params.ID,
			Email:        params.Email,
			PasswordHash: params.PasswordHash,
			CreatedAt:    params.CreatedAt,
			UpdatedAt:    params.UpdatedAt,
		},
	)

	var postgresError *pq.Error

	if errors.As(err, &postgresError) &&
		postgresError.Code == "23505" &&
		postgresError.Constraint == "accounts_email_unique" {
		return account.ErrEmailTaken
	}

	return err
}

func (store *AccountStore) FindByEmail(
	ctx context.Context,
	email string,
) (account.Account, error) {
	acc, err := store.queries.FindAccountByEmail(
		ctx,
		email,
	)

	if err != nil {
		return account.Account{}, err
	}

	return account.Account{
		ID:           acc.ID,
		PasswordHash: acc.PasswordHash,
		Email:        acc.Email,
		CreatedAt:    acc.CreatedAt,
		UpdatedAt:    acc.UpdatedAt,
	}, nil
}

var _ account.Store = (*AccountStore)(nil)
