package postgres

import (
	"context"
	"database/sql"

	"github.com/getshelf/backend/internal/adapters/postgres/sqlcgen"
	c "github.com/getshelf/backend/internal/modules/collection"
)

type CollectionStore struct {
	queries *sqlcgen.Queries
}

func NewCollectionStore(db *sql.DB) *CollectionStore {
	return &CollectionStore{
		queries: sqlcgen.New(db),
	}
}

func nullableString(value *string) sql.NullString {
	if value == nil {
		return sql.NullString{}
	}

	return sql.NullString{
		String: *value,
		Valid:  true,
	}
}

func (store *CollectionStore) Create(
	ctx context.Context,
	params c.CreateCollectionParams,
) (c.Collection, error) {
	collection, err := store.queries.CreateCollection(
		ctx,
		sqlcgen.CreateCollectionParams{
			ID:        params.ID,
			Title:     params.Title,
			Icon:      nullableString(params.Icon),
			ParentID:  nullableString(params.ParentID),
			OwnerID:   params.OwnerID,
			CreatedAt: params.CreatedAt,
			UpdatedAt: params.UpdatedAt,
		},
	)

	if err != nil {
		return c.Collection{}, err
	}

	var parentID *string
	if collection.ParentID.Valid {
    	parentID = &collection.ParentID.String
	}

	return c.Collection{
		ID:        collection.ID,
		Title:     collection.Title,
		Icon:      &collection.Icon.String,
		ParentID:  parentID,
		OwnerID:   collection.OwnerID,
		SortOrder: int(collection.SortOrder),
		CreatedAt: collection.CreatedAt,
		UpdatedAt: collection.UpdatedAt,
	}, nil
}
