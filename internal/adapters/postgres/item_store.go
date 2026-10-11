package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/getshelf/backend/internal/adapters/postgres/sqlcgen"
	"github.com/getshelf/backend/internal/modules/item"
)

type ItemStore struct {
	db *sql.DB
	queries *sqlcgen.Queries
}

func NewItemStore(db *sql.DB) *ItemStore {
	return &ItemStore{
		db: db,
		queries: sqlcgen.New(db),
	}
}

var _ item.Store = (*ItemStore)(nil)

func (store *ItemStore) SaveURL(
	ctx context.Context,
	params item.SaveURLParams,
) (item.Item, error) {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return item.Item{}, fmt.Errorf("Begin save item: %w", err)
	}
	defer tx.Rollback()

	queries := store.queries.WithTx(tx)

	row, err := queries.CreateSavedURL(ctx, sqlcgen.CreateSavedURLParams{
		ID: params.ID,
		OwnerID: params.OwnerID,
		SourceUrl: sql.NullString{
			String: params.URL,
			Valid: true,
		},
		CreatedAt: params.CreatedAt,
		UpdatedAt: params.CreatedAt,
	})

	if err != nil {
		return item.Item{}, fmt.Errorf("Create extraction: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return item.Item{}, fmt.Errorf("Commit save item: %w", err)
	}

	return item.Item{
		ID: row.ID,
		OwnerID: row.OwnerID,
		Kind: row.Kind,
		Title: stringPointer(row.Title),
		SourceURL: stringPointer(row.SourceUrl),
		ContentText: stringPointer(row.ContentText),
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}

func stringPointer(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}

	return &value.String
}
