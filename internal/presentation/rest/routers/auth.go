package routers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	applicationUser "github.com/getshelf/backend/internal/application/user"
	domainErrors "github.com/getshelf/backend/internal/domain/user/errors"
	"github.com/getshelf/backend/internal/domain/user/ports"
	"github.com/getshelf/backend/internal/domain/user/values"
	"github.com/getshelf/backend/internal/presentation/rest/schemas"
)

type AuthHandler struct {
	service applicationUser.AuthService
}

func NewAuthHandler(service applicationUser.AuthService) AuthHandler {
	return AuthHandler{service: service}
}

func (handler AuthHandler) Login(responseWriter http.ResponseWriter, request *http.Request) {
	var input schemas.LoginRequest
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		writeError(responseWriter, http.StatusBadRequest, "invalid JSON body")
		return
	}
	pair, err := handler.service.Login(request.Context(), applicationUser.LoginRequest{Login: input.Login, Password: input.Password})
	if err != nil {
		writeAuthError(responseWriter, err)
		return
	}
	writeJSON(responseWriter, http.StatusOK, schemas.TokenResponse{AccessToken: pair.AccessToken, RefreshToken: pair.RefreshToken})
}

func (handler AuthHandler) Refresh(responseWriter http.ResponseWriter, request *http.Request) {
	var input schemas.RefreshRequest
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		writeError(responseWriter, http.StatusBadRequest, "invalid JSON body")
		return
	}
	pair, err := handler.service.Refresh(request.Context(), input.RefreshToken)
	if err != nil {
		writeAuthError(responseWriter, err)
		return
	}
	writeJSON(responseWriter, http.StatusOK, schemas.TokenResponse{AccessToken: pair.AccessToken, RefreshToken: pair.RefreshToken})
}

func (handler AuthHandler) Logout(responseWriter http.ResponseWriter, request *http.Request) {
	access := bearerToken(request)
	var input schemas.RefreshRequest
	_ = json.NewDecoder(request.Body).Decode(&input)
	if err := handler.service.Logout(request.Context(), access, input.RefreshToken); err != nil {
		writeError(responseWriter, http.StatusInternalServerError, "could not revoke tokens")
		return
	}
	responseWriter.WriteHeader(http.StatusNoContent)
}

func (handler AuthHandler) Profile(responseWriter http.ResponseWriter, request *http.Request) {
	profile, err := handler.service.Profile(request.Context(), UserID(request.Context()))
	if err != nil {
		writeAuthError(responseWriter, err)
		return
	}
	writeJSON(responseWriter, http.StatusOK, profileResponse(profile))
}

func (handler AuthHandler) PatchProfile(responseWriter http.ResponseWriter, request *http.Request) {
	var input schemas.UpdateProfileRequest
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		writeError(responseWriter, http.StatusBadRequest, "invalid JSON body")
		return
	}
	profile, err := handler.service.UpdateProfile(request.Context(), UserID(request.Context()), applicationUser.UpdateProfileRequest{Username: input.Username, Email: input.Email})
	if err != nil {
		writeAuthError(responseWriter, err)
		return
	}
	writeJSON(responseWriter, http.StatusOK, profileResponse(profile))
}

func profileResponse(profile applicationUser.Profile) schemas.ProfileResponse {
	return schemas.ProfileResponse{ID: profile.ID.String(), Username: profile.Username.String(), Email: profile.Email.String(), CreatedAt: profile.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00")}
}

func writeAuthError(responseWriter http.ResponseWriter, err error) {
	status := http.StatusUnauthorized
	switch {
	case errors.Is(err, domainErrors.ErrUsernameTaken), errors.Is(err, domainErrors.ErrEmailTaken):
		status = http.StatusConflict
	case errors.Is(err, domainErrors.ErrInvalidUsername), errors.Is(err, domainErrors.ErrInvalidEmail):
		status = http.StatusBadRequest
	case errors.Is(err, domainErrors.ErrUserNotFound):
		status = http.StatusNotFound
	}
	writeError(responseWriter, status, err.Error())
}

func bearerToken(request *http.Request) string {
	parts := strings.Fields(request.Header.Get("Authorization"))
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}
	return ""
}

func UserID(ctx context.Context) values.UserID {
	id, _ := ctx.Value(userIDKey{}).(values.UserID)
	return id
}

type userIDKey struct{}

func withUserID(ctx context.Context, id values.UserID) context.Context {
	return context.WithValue(ctx, userIDKey{}, id)
}

func NewAuthMiddleware(tokens ports.TokenService, blacklist ports.TokenBlacklist) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
			token := bearerToken(request)
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
			next.ServeHTTP(responseWriter, request.WithContext(withUserID(request.Context(), claims.UserID)))
		})
	}
}

func writeJSON(responseWriter http.ResponseWriter, status int, body any) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(status)
	_ = json.NewEncoder(responseWriter).Encode(body)
}
