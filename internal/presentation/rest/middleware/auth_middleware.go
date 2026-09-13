package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	domainErrors "github.com/getshelf/backend/internal/domain/user/errors"
	"github.com/getshelf/backend/internal/domain/user/ports"
	"github.com/getshelf/backend/internal/domain/user/values"
	"github.com/getshelf/backend/internal/presentation/rest/schemas"
)

func BearerToken(request *http.Request) string {
	parts := strings.Fields(request.Header.Get("Authorization"))
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}
	return ""
}

func WriteAuthError(responseWriter http.ResponseWriter, err error) {
	status := http.StatusUnauthorized
	switch {
	case errors.Is(err, domainErrors.ErrEmailTaken):
		status = http.StatusConflict
	case errors.Is(err, domainErrors.ErrInvalidEmail):
		status = http.StatusBadRequest
	case errors.Is(err, domainErrors.ErrUserNotFound):
		status = http.StatusNotFound
	}
	writeError(responseWriter, status, err.Error())
}

func NewAuthMiddleware(tokens ports.TokenService, blacklist ports.TokenBlacklist) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
			token := BearerToken(request)
			claims, err := tokens.Parse(token, false)
			if err != nil {
				writeError(responseWriter, http.StatusUnauthorized, domainErrors.ErrUnauthorized.Error())
				return
			}
			revoked, err := blacklist.IsRevoked(request.Context(), claims.TokenID)
			if err != nil || revoked {
				writeError(responseWriter, http.StatusUnauthorized, domainErrors.ErrUnauthorized.Error())
				return
			}
			next.ServeHTTP(responseWriter, request.WithContext(WithUserID(request.Context(), claims.UserID)))
		})
	}
}

func UserID(ctx context.Context) values.UserID {
	id, _ := ctx.Value(userIDKey{}).(values.UserID)
	return id
}

func WithUserID(ctx context.Context, id values.UserID) context.Context {
	return context.WithValue(ctx, userIDKey{}, id)
}

type userIDKey struct{}

func writeError(responseWriter http.ResponseWriter, status int, message string) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(status)
	_ = json.NewEncoder(responseWriter).Encode(schemas.ErrorResponse{Error: message})
}
