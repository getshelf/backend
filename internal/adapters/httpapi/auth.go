package httpapi

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/getshelf/backend/internal/modules/account"
	"github.com/getshelf/backend/internal/modules/session"
)

type loginInput struct {
	Body struct {
		Email    string `json:"email" format:"email"`
		Password string `json:"password"`
	}
}

type loginOutput struct {
	SetCookie http.Cookie `header:"Set-Cookie"`

	Body struct {
		Account accountResponse `json:"account"`
	}
}

type accountResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func registerAuthRoutes(
	api huma.API,
	accounts *account.Service,
	sessions *session.Service,
	secureCookies bool,
	logger *slog.Logger,
) {
	huma.Register(
		api,
		huma.Operation{
			OperationID: "login",
			Method: http.MethodPost,
			Path: "/auth/login",
			Summary: "Log in",
			Tags: []string{"Authentication"},

			Errors: []int{
				http.StatusUnauthorized,
			},
		},
		func(
			ctx context.Context,
			input *loginInput,
		) (*loginOutput, error) {
			authenticated, err := accounts.Authenticate(
				ctx,
				account.AuthenticateInput{
					Email:    input.Body.Email,
					Password: input.Body.Password,
				},
			)
			if err != nil {
				return nil, huma.Error401Unauthorized(
					"invalid email or password",
				)
			}

			createdSession, err := sessions.Create(
				ctx,
				authenticated.ID,
			)
			if err != nil {
				logger.Error("create session", "error", err)
				return nil, huma.Error500InternalServerError(
					"could not create session",
				)
			}

			output := &loginOutput{
				SetCookie: newSessionCookie(
					createdSession.Token,
					createdSession.ExpiresAt,
					secureCookies,
				),
			}

			output.Body.Account = accountResponse{
				ID:    authenticated.ID,
				Email: authenticated.Email,
			}

			return output, nil
		},
	)
}
