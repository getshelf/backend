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

type logoutInput struct {
    Session string `cookie:"session" doc:"Session token set by the login endpoint."`
}
type logoutOutput struct {
    SetCookie http.Cookie `header:"Set-Cookie" doc:"Expired session cookie that instructs the browser to delete it."`
}

func registerAuthRoutes(
	api huma.API,
	accounts *account.Service,
	sessions *session.Service,
	secureCookies bool,
	logger *slog.Logger,
) {
	// Login
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
					session.SessionCookieName,
				),
			}

			output.Body.Account = accountResponse{
				ID:    authenticated.ID,
				Email: authenticated.Email,
			}

			return output, nil
		},
	)

	// Logout
	huma.Register(
		api,
		huma.Operation{
			OperationID: "logout",
			Method: http.MethodPost,
			Path: "/auth/logout",
			Summary: "Log out",
			Description: "Ends the current session and clears the session cookie",
			Tags: []string{"Authentication"},
			DefaultStatus: http.StatusNoContent,
		},
		func(
			ctx context.Context,
			input *logoutInput,
		) (*logoutOutput, error) {
			if err := sessions.Logout(ctx, input.Session); err != nil {
				logger.Error("logout", "error", err)
				return nil, huma.Error500InternalServerError("could not log out")
			}

			return &logoutOutput{
				SetCookie: clearSessionCookie(secureCookies, session.SessionCookieName),
			}, nil
		},
	)
}
