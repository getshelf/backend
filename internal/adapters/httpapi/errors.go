package httpapi
import (
	"context"
	"errors"
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
	"github.com/getshelf/backend/internal/modules/account"
)

func mapAccountError(
	ctx context.Context,
	logger *slog.Logger,
	err error,
) error {
	switch {
		case errors.Is(err, account.ErrInvalidEmail):
			return huma.Error422UnprocessableEntity(
				"invalid email",
			)

		case errors.Is(err, account.ErrPasswordTooWeak):
			return huma.Error422UnprocessableEntity(
				"password must contain at least eight characters",
			)

		case errors.Is(err, account.ErrPasswordContainsEmail):
			return huma.Error422UnprocessableEntity(
				"password must not contain the email",
			)

		case errors.Is(err, account.ErrEmailTaken):
			return huma.Error409Conflict(
				"email is already taken",
			)

		default:
			logger.ErrorContext(
				ctx,
				"account operation failed",
				"error",
				err,
			)

			return huma.Error500InternalServerError(
				"unexpected server error",
			)
		}
}
