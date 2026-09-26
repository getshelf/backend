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

func (store *CollectionStore) Update(
	ctx context.Context,
	collectionID string,
	ownerID string,
	params c.UpdateCollectionParams,
) (c.Collection, error) {
	collection, err := store.queries.UpdateCollection(
		ctx,
		sqlcgen.UpdateCollectionParams{
			ID:        collectionID,
			Title:     params.Title,
			Icon:      nullableString(params.Icon),
			ParentID:  nullableString(params.ParentID.Value),
			OwnerID:   ownerID,
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
		SortOrder: int(collection.SortOrder),
		CreatedAt: collection.CreatedAt,
		UpdatedAt: collection.UpdatedAt,
	}, nil
}

func (store *CollectionStore) Get(
	ctx context.Context,
	collectionID string,
	ownerID string,
) (c.Collection, error) {
	collection, err := store.queries.GetCollection(
		ctx,
		sqlcgen.GetCollectionParams{
			ID:      collectionID,
			OwnerID: ownerID,
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
		SortOrder: int(collection.SortOrder),
		CreatedAt: collection.CreatedAt,
		UpdatedAt: collection.UpdatedAt,
	}, nil
}

func (store *CollectionStore) List(
	ctx context.Context,
	ownerID *string,
) ([]c.Collection, error) {
	collections, err := store.queries.ListCollections(
		ctx,
		*ownerID,
	)

	if err != nil {
		return nil, err
	}

	result := nest(collections)

	return result, nil
}

func nest(collections []sqlcgen.Collection) []c.Collection {
	collectionsOf := make(map[string][]c.Collection, len(collections))
	for _, collection := range collections {
		key := ""
		if collection.ParentID.Valid {
			key = *&collection.ParentID.String
		}
		collectionsOf[key] = append(collectionsOf[key], c.Collection{
			ID:        collection.ID,
			Title:     collection.Title,
			Icon:      &collection.Icon.String,
			ParentID:  &collection.ParentID.String,
			SortOrder: int(collection.SortOrder),
			CreatedAt: collection.CreatedAt,
			UpdatedAt: collection.UpdatedAt,
			Collections: []c.Collection{},
		})
	}

	var walk func(parentKey string) []c.Collection
	walk = func(parentKey string) []c.Collection {
		collections := collectionsOf[parentKey]
		items := make([]c.Collection, len(collections))

		for i, collection := range collections {
			items[i] = c.Collection{
				ID:        collection.ID,
				Title:     collection.Title,
				Icon:      collection.Icon,
				ParentID:  collection.ParentID,
				SortOrder: collection.SortOrder,
				CreatedAt: collection.CreatedAt,
				UpdatedAt: collection.UpdatedAt,
				Collections: walk(collection.ID),
			}
		}
		return items
	}

	return walk("")
}
