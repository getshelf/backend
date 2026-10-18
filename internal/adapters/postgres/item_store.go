package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/getshelf/backend/internal/adapters/postgres/sqlcgen"
	"github.com/getshelf/backend/internal/modules/collection"
	"github.com/getshelf/backend/internal/modules/item"
)

type ItemStore struct {
	db      *sql.DB
	queries *sqlcgen.Queries
}

type collectionJSON struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Icon      *string   `json:"icon"`
	ParentID  *string   `json:"parent_id"`
	OwnerID   string    `json:"owner_id"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewItemStore(db *sql.DB) *ItemStore {
	return &ItemStore{
		db:      db,
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
		ID:      params.ID,
		OwnerID: params.OwnerID,
		SourceUrl: sql.NullString{
			String: params.URL,
			Valid:  true,
		},
		CollectionID: nullableString(params.CollectionID),
		CreatedAt:    params.CreatedAt,
		UpdatedAt:    params.CreatedAt,
	})

	if err != nil {
		return item.Item{}, fmt.Errorf("Create extraction: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return item.Item{}, fmt.Errorf("Commit save item: %w", err)
	}

	savedCollection, err := decodeCollection(row.Collection.([]byte))

	if err != nil {
		return item.Item{}, err
	}

	return item.Item{
		ID:          row.ID,
		OwnerID:     row.OwnerID,
		Kind:        row.Kind,
		Title:       stringPointer(row.Title),
		SourceURL:   stringPointer(row.SourceUrl),
		ContentText: stringPointer(row.ContentText),
		Collection:  savedCollection,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}, nil
}

func stringPointer(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}

	return &value.String
}

func decodeCollection(data []byte) (*collection.Collection, error) {
	if len(data) == 0 {
		return nil, nil
	}

	var value *collectionJSON

	if err := json.Unmarshal(data, &value); err != nil {
		return nil, fmt.Errorf("decode item collection: %w", err)
	}

	if value == nil {
		return nil, nil
	}

	return &collection.Collection{
		ID:        value.ID,
		Title:     value.Title,
		Icon:      value.Icon,
		ParentID:  value.ParentID,
		OwnerID:   value.OwnerID,
		SortOrder: value.SortOrder,
		CreatedAt: value.CreatedAt,
		UpdatedAt: value.UpdatedAt,
	}, nil
}
