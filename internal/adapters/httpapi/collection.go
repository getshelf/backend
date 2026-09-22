package httpapi

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/getshelf/backend/internal/modules/collection"
	"github.com/getshelf/backend/internal/modules/session"
)

type createCollectionInput struct {
	Body createCollectionInputBody
}

type createCollectionInputBody struct {
	Title    string  `json:"title" binding:"required"`
	Icon     *string `json:"icon"`
	ParentID *string `json:"parent_id" nullable:"true"`
}

type createCollectionOutput struct {
	Body createCollectionOutputBody
}

type createCollectionOutputBody struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Icon      *string   `json:"icon"`
	ParentID  *string   `json:"parent_id" nullable:"true"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
			OperationID:   "createCollection",
			Method:        http.MethodPost,
			Path:          "/collections",
			Summary:       "Create a Shelf collection",
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

			fmt.Println("collection:AccountID", accountID)

			if !ok {
				return nil, huma.Error401Unauthorized("unauthorized")
			}


			created, err := collections.CreateCollection(
				ctx,
				collection.CreateCollectionInput{
					Title:     input.Body.Title,
					Icon:      input.Body.Icon,
					ParentID:  nil,
					OwnerID:   accountID,
				},
			)

			if err != nil {
				return nil, err
			}

			return &createCollectionOutput{
				Body: createCollectionOutputBody{
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
}
