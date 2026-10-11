package httpapi

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/getshelf/backend/internal/modules/item"
	"github.com/getshelf/backend/internal/modules/session"
)

type Item struct {
	ID          string    `json:"id"`
	Kind        string    `json:"kind"`
	Title       *string   `json:"title"`
	SourceURL   *string   `json:"sourceUrl"`
	ContentText *string   `json:"contentText"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type createItemInput struct {
	Body createItemInputBody
}

type createItemInputBody struct {
	URL string `json:"url"`
}

type createItemOutput struct {
	Body Item
}

func registerItemRoutes(
	api huma.API,
	items *item.Service,
	session *session.Service,
	_ *slog.Logger,
) {
	huma.Register(
		api,
		huma.Operation{
			OperationID:   "createItem",
			Method:        http.MethodPost,
			Path:          "/items",
			Summary:       "Create an item",
			Tags:          []string{"Items"},
			DefaultStatus: http.StatusCreated,
		},
		func(
			ctx context.Context,
			input *createItemInput,
		) (*createItemOutput, error) {
			accountID, ok := session.GetAccountID(ctx)

			if !ok {
				return nil, huma.Error401Unauthorized("unauthorized")
			}

			item, err := items.SaveURL(ctx, accountID, input.Body.URL)

			if err != nil {
				return nil, err
			}

			return &createItemOutput{
				Body: Item{
					ID:          item.ID,
					Title:       item.Title,
					Kind:        "link",
					SourceURL:   item.SourceURL,
					ContentText: item.ContentText,
					CreatedAt:   item.CreatedAt,
					UpdatedAt:   item.UpdatedAt,
				},
			}, nil
		},
	)
}
