package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"reflect"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/getshelf/backend/internal/modules/collection"
	"github.com/getshelf/backend/internal/modules/session"
)

type OmittableNullable[T any] struct {
	Sent  bool
	Null  bool
	Value T
}

type CollectionItem struct {
	ID          string           `json:"id"`
	Title       string           `json:"title"`
	Icon        *string          `json:"icon"`
	ParentID    *string          `json:"parent_id"`
	SortOrder   int              `json:"sort_order"`
	Collections []CollectionItem `json:"collections"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type getCollectionInput struct {
	CollectionID string `path:"collectionID"`
}

type getCollectionOutput struct {
	Body CollectionItem
}

type getCollectionsInput struct {
}

type getCollectionsOutput struct {
	Body getCollectionsOutputBody
}

type getCollectionsOutputBody struct {
	Collections []CollectionItem `json:"collections"`
}

type createCollectionInput struct {
	Body createCollectionInputBody
}

type createCollectionInputBody struct {
	Title    string  `json:"title"`
	Icon     *string `json:"icon"`
	ParentID *string `json:"parent_id"`
}

type createCollectionOutput struct {
	Body CollectionItem
}

type updateCollectionInput struct {
	CollectionID string `path:"collectionID"`
	Body         updateCollectionInputBody
}

type updateCollectionInputBody struct {
	Title    string                    `json:"title"`
	Icon     *string                   `json:"icon"`
	ParentID OmittableNullable[string] `json:"parent_id,omitempty" required:"false"`
}

type updateCollectionOutput struct {
	Body CollectionItem
}

func registerCollectionRoutes(
	api huma.API,
	collections *collection.Service,
	session *session.Service,
	_ *slog.Logger,
) {
	huma.Register(
		api,
		huma.Operation{
			OperationID:   "getCollection",
			Method:        http.MethodGet,
			Path:          "/collections/{collectionID}",
			Summary:       "Get a collection",
			Tags:          []string{"Collections"},
			DefaultStatus: http.StatusOK,
		},
		func(
			ctx context.Context,
			input *getCollectionInput,
		) (*getCollectionOutput, error) {
			accountID, ok := session.GetAccountID(ctx)

			if !ok {
				return nil, huma.Error401Unauthorized("unauthorized")
			}

			collection, err := collections.GetCollection(ctx, input.CollectionID, accountID)
			if err != nil {
				return nil, err
			}

			return &getCollectionOutput{
				Body: CollectionItem{
					ID:        collection.ID,
					Title:     collection.Title,
					Icon:      collection.Icon,
					ParentID:  collection.ParentID,
					SortOrder: collection.SortOrder,
					CreatedAt: collection.CreatedAt,
					UpdatedAt: collection.UpdatedAt,
				},
			}, nil
		},
	)

	huma.Register(
		api,
		huma.Operation{
			OperationID:   "getCollections",
			Method:        http.MethodGet,
			Path:          "/collections",
			Summary:       "Get all collections",
			Tags:          []string{"Collections"},
			DefaultStatus: http.StatusOK,
		},
		func(
			ctx context.Context,
			input *getCollectionsInput,
		) (*getCollectionsOutput, error) {
			accountID, ok := session.GetAccountID(ctx)

			if !ok {
				return nil, huma.Error401Unauthorized("unauthorized")
			}

			collections, err := collections.ListCollections(ctx, accountID)
			if err != nil {
				return nil, err
			}

			items := make([]CollectionItem, 0, len(collections))
			for _, collection := range collections {
				items = append(items, CollectionItem{
					ID:          collection.ID,
					Title:       collection.Title,
					Icon:        collection.Icon,
					ParentID:    collection.ParentID,
					SortOrder:   collection.SortOrder,
					Collections: nestCollections(collection.Collections),
					CreatedAt:   collection.CreatedAt,
					UpdatedAt:   collection.UpdatedAt,
				})
			}

			return &getCollectionsOutput{
				Body: getCollectionsOutputBody{
					Collections: items,
				},
			}, nil
		},
	)

	huma.Register(
		api,
		huma.Operation{
			OperationID:   "createCollection",
			Method:        http.MethodPost,
			Path:          "/collections",
			Summary:       "Create a collection",
			Tags:          []string{"Collections"},
			DefaultStatus: http.StatusCreated,

			Errors: []int{
				http.StatusConflict,
			},
		},
		func(
			ctx context.Context,
			input *createCollectionInput,
		) (*createCollectionOutput, error) {
			accountID, ok := session.GetAccountID(ctx)

			if !ok {
				return nil, huma.Error401Unauthorized("unauthorized")
			}

			created, err := collections.CreateCollection(
				ctx,
				collection.CreateCollectionInput{
					Title:    input.Body.Title,
					Icon:     input.Body.Icon,
					ParentID: input.Body.ParentID,
					OwnerID:  accountID,
				},
			)

			if err != nil {
				return nil, err
			}

			return &createCollectionOutput{
				Body: CollectionItem{
					ID:        created.ID,
					Title:     created.Title,
					Icon:      created.Icon,
					ParentID:  created.ParentID,
					SortOrder: created.SortOrder,
					CreatedAt: created.CreatedAt,
					UpdatedAt: created.UpdatedAt,
				},
			}, nil
		},
	)

	huma.Register(
		api,
		huma.Operation{
			OperationID:   "patchCollection",
			Method:        http.MethodPatch,
			Path:          "/collections/{collectionID}",
			Summary:       "Patch a collection",
			Tags:          []string{"Collections"},
			DefaultStatus: http.StatusOK,
		},
		func(
			ctx context.Context,
			input *updateCollectionInput,
		) (*updateCollectionOutput, error) {
			accountID, ok := session.GetAccountID(ctx)

			if !ok {
				return nil, huma.Error401Unauthorized("unauthorized")
			}

			patched, err := collections.UpdateCollection(
				ctx,
				input.CollectionID,
				accountID,
				collection.UpdateCollectionParams{
					Title: input.Body.Title,
					Icon:  input.Body.Icon,
					ParentID: collection.ParentID{
						Set:   input.Body.ParentID.Sent,
						Value: &input.Body.ParentID.Value,
					},
				},
			)

			if err != nil {
				return nil, err
			}

			return &updateCollectionOutput{
				Body: CollectionItem{
					ID:        patched.ID,
					Title:     patched.Title,
					Icon:      patched.Icon,
					ParentID:  patched.ParentID,
					SortOrder: patched.SortOrder,
					CreatedAt: patched.CreatedAt,
					UpdatedAt: patched.UpdatedAt,
				},
			}, nil
		},
	)
}

func nestCollections(collections []collection.Collection) []CollectionItem {
	items := make([]CollectionItem, 0, len(collections))
	for _, collection := range collections {
		items = append(items, CollectionItem{
			ID:          collection.ID,
			Title:       collection.Title,
			Icon:        collection.Icon,
			ParentID:    collection.ParentID,
			SortOrder:   collection.SortOrder,
			CreatedAt:   collection.CreatedAt,
			UpdatedAt:   collection.UpdatedAt,
			Collections: nestCollections(collection.Collections),
		})
	}
	return items
}

func (o *OmittableNullable[T]) UnmarshalJSON(b []byte) error {
	o.Sent = true

	if bytes.Equal(b, []byte("null")) {
		o.Null = true
		return nil
	}

	return json.Unmarshal(b, &o.Value)
}

func (o OmittableNullable[T]) Schema(r huma.Registry) *huma.Schema {
	schema := r.Schema(reflect.TypeOf(o.Value), true, "")

	schema.Nullable = true

	return schema
}
