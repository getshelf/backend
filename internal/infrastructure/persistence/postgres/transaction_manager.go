package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type TransactionManager struct {
	db *sqlx.DB
}

func NewTransactionManager(db *sqlx.DB) TransactionManager {
	return TransactionManager{db: db}
}

func (manager TransactionManager) Do(ctx context.Context, callback func(context.Context) error) error {
	tx, err := manager.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	transactionContext := ContextWithTx(ctx, tx)
	if err := callback(transactionContext); err != nil {
		return err
	}
	return tx.Commit()
}

type transactionKey struct{}

func ContextWithTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, transactionKey{}, tx)
}

func TxFromContext(ctx context.Context) (*sqlx.Tx, bool) {
	tx, ok := ctx.Value(transactionKey{}).(*sqlx.Tx)
	return tx, ok
}
