package integration

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	applicationUser "github.com/getshelf/backend/internal/application/user"
	"github.com/getshelf/backend/internal/domain/user/ports"
	"github.com/getshelf/backend/internal/domain/user/values"
	"github.com/getshelf/backend/internal/presentation/rest/routers"
)

type authHTTPRepository struct {
	user ports.AuthUser
}

func (fake authHTTPRepository) FindByLogin(context.Context, string) (ports.AuthUser, error) {
	return fake.user, nil
}

func (fake authHTTPRepository) FindByID(context.Context, values.UserID) (ports.AuthUser, error) {
	return fake.user, nil
}

func (fake authHTTPRepository) UpdateProfile(context.Context, values.UserID, values.Username, values.Email) error {
	return nil
}

type authHTTPVerifier struct{}

func (authHTTPVerifier) Compare(values.PasswordHash, string) error { return nil }

type authHTTPTokenService struct {
	pair ports.TokenPair
}

func (fake authHTTPTokenService) Issue(ports.AuthUser) (ports.TokenPair, error) {
	return fake.pair, nil
}

func (authHTTPTokenService) Parse(raw string, refresh bool) (ports.TokenClaims, error) {
	if raw == "" {
		return ports.TokenClaims{}, errors.New("missing token")
	}
	tokenID := raw
	if refresh {
		tokenID = "refresh-id"
	}
	return ports.TokenClaims{UserID: values.NewUserID("user-1"), TokenID: tokenID, ExpiresAt: time.Now().Add(time.Hour), Refresh: refresh}, nil
}

type authHTTPBlacklist struct {
	revoked map[string]bool
}

func (fake *authHTTPBlacklist) Revoke(_ context.Context, tokenID string, _ time.Duration) error {
	fake.revoked[tokenID] = true
	return nil
}

func (fake *authHTTPBlacklist) IsRevoked(_ context.Context, tokenID string) (bool, error) {
	return fake.revoked[tokenID], nil
}

type authHTTPChecker struct{}

func (authHTTPChecker) IsUsernameTaken(context.Context, values.Username) (bool, error) {
	return false, nil
}
func (authHTTPChecker) IsEmailTaken(context.Context, values.Email) (bool, error) { return false, nil }

func newAuthRouter() http.Handler {
	username, _ := values.NewUsername("john_doe")
	email, _ := values.NewEmail("john@example.com")
	hash, _ := values.NewPasswordHash("hash")
	service := applicationUser.NewAuthService(
		authHTTPRepository{user: ports.AuthUser{ID: values.NewUserID("user-1"), Username: username, Email: email, PasswordHash: hash, CreatedAt: time.Unix(0, 0).UTC()}},
		authHTTPVerifier{},
		authHTTPTokenService{pair: ports.TokenPair{AccessToken: "access-token", RefreshToken: "refresh-token"}},
		&authHTTPBlacklist{revoked: map[string]bool{}},
		authHTTPChecker{},
	)
	return routers.NewRouter(
		routers.RegisterHandler{},
		routers.NewAuthHandler(service),
		func(next http.Handler) http.Handler {
			return routers.NewAuthMiddleware(authHTTPTokenService{}, &authHTTPBlacklist{revoked: map[string]bool{}})(next)
		},
	)
}

func TestAuthHTTPHandlers(t *testing.T) {
	router := newAuthRouter()

	tests := []struct {
		name   string
		method string
		path   string
		body   string
		token  string
		want   int
		match  string
	}{
		{name: "login", method: http.MethodPost, path: "/api/v1/users/login", body: `{"login":"john@example.com","password":"StrongPass123"}`, want: http.StatusOK, match: `"access_token":"access-token"`},
		{name: "refresh", method: http.MethodPost, path: "/api/v1/users/refresh", body: `{"refresh_token":"refresh-token"}`, want: http.StatusOK, match: `"refresh_token":"refresh-token"`},
		{name: "profile get", method: http.MethodGet, path: "/api/v1/users/profile", token: "access-token", want: http.StatusOK, match: `"username":"john_doe"`},
		{name: "profile patch", method: http.MethodPatch, path: "/api/v1/users/profile", body: `{"username":"jane_doe","email":"jane@example.com"}`, token: "access-token", want: http.StatusOK, match: `"email":"jane@example.com"`},
		{name: "logout", method: http.MethodPost, path: "/api/v1/users/logout", body: `{"refresh_token":"refresh-token"}`, token: "access-token", want: http.StatusNoContent},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.path, strings.NewReader(test.body))
			if test.token != "" {
				request.Header.Set("Authorization", "Bearer "+test.token)
			}
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			assert.Equal(t, test.want, response.Code, "body: %s", response.Body.String())
			if test.match != "" && !strings.Contains(response.Body.String(), test.match) {
				assert.Contains(t, response.Body.String(), test.match)
			}
		})
	}
}

func TestAuthMiddlewareRejectsMissingBearerToken(t *testing.T) {
	router := newAuthRouter()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/users/profile", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	require.Equal(t, http.StatusUnauthorized, response.Code)
}
