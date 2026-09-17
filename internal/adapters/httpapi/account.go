package httpapi

import (
	"context"
	"net/http"
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
	"github.com/getshelf/backend/internal/modules/account"
)

type registerAccountInput struct {
	Body registerAccountInputBody
}

type registerAccountInputBody struct {
	Email    string `json:"email" format:"email" doc:"The email address of the account"`
	Password string `json:"password" doc:"The password of the account"`
}

type registerAccountOutput struct {
	Body registerAccountResponseBody
}

type registerAccountResponseBody struct {
	ID       string `json:"id" doc:"The ID of the account"`
}

func registerAccountRoutes(
	api huma.API,
	accounts *account.Service,
	logger *slog.Logger,
) {
	huma.Register(
		api,
		huma.Operation{
			OperationID:   "registerAccount",
			Method:        http.MethodPost,
			Path:          "/users/register",
			Summary:       "Register an account",
			Description:   "Create a Shelf account using an email and password.",
			Tags:          []string{"Users"},
			DefaultStatus: http.StatusCreated,

			// Huma will also document validation and internal errors.
			Errors: []int{
				http.StatusConflict,
			},
		},
		func(
			ctx context.Context,
			input *registerAccountInput,
		) (*registerAccountOutput, error) {
			created, err := accounts.Register(
				ctx,
				account.RegisterInput{
					Email:    input.Body.Email,
					Password: input.Body.Password,
				},
			)

			if err != nil {
				return nil, mapAccountError(
					ctx,
					logger,
					err,
				)
			}

			return &registerAccountOutput{
				Body: registerAccountResponseBody{
					ID: created.ID,
				},
			}, nil
		},
	)
}
