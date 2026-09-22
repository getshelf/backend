package middlewares

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/getshelf/backend/internal/modules/session"
)


func RequireSession(
	api huma.API,
	sessions *session.Service,
	sessionCookieName string,
) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		cookie, err := huma.ReadCookie(ctx, sessionCookieName)
		if err != nil || cookie == nil || cookie.Value == "" {
			_ = huma.WriteErr(
				api,
				ctx,
				http.StatusUnauthorized,
				"authentication required",
			)
			return
		}

		accountID, err := sessions.Authenticate(
			ctx.Context(),
			cookie.Value,
		)

		if err != nil {
			_ = huma.WriteErr(
				api,
				ctx,
				http.StatusUnauthorized,
				"invalid or expired session",
			)
			return
		}

		next(huma.WithValue(
			ctx,
			session.AccountIDContextKey{},
			accountID,
		))
	}
}
