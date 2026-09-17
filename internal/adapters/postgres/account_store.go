package postgres

import (
	"context"
	"database/sql"

	"github.com/getshelf/backend/internal/modules/account"
	// "github.com/lib/pq"
)

type AccountStore struct {
	db *sql.DB
}

func NewAccountStore(db *sql.DB) *AccountStore {
	return &AccountStore{db: db}
}

func (store *AccountStore) Create (
	ctx context.Context,
	params account.CreateAccountParams,
) error {
	return nil
}

func (store *AccountStore) FindByEmail(
	ctx context.Context,
	email string,
) (string, error) {

	// TODO: implement
	return "", nil
}

var _ account.Store = (*AccountStore)(nil)
