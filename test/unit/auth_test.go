package unit

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	applicationUser "github.com/getshelf/backend/internal/application/user"
	domainErrors "github.com/getshelf/backend/internal/domain/user/errors"
	"github.com/getshelf/backend/internal/domain/user/ports"
	"github.com/getshelf/backend/internal/domain/user/values"
)

type authRepositoryFake struct {
	user         ports.AuthUser
	findErr      error
	updatedID    values.UserID
	updatedEmail values.Email
}

func (fake *authRepositoryFake) FindByLogin(context.Context, string) (ports.AuthUser, error) {
	return fake.user, fake.findErr
}

func (fake *authRepositoryFake) FindByID(context.Context, values.UserID) (ports.AuthUser, error) {
	return fake.user, fake.findErr
}

func (fake *authRepositoryFake) UpdateProfile(_ context.Context, id values.UserID, email values.Email) error {
	fake.updatedID, fake.updatedEmail = id, email
	return nil
}

type verifierFake struct{ err error }

func (fake verifierFake) Compare(values.PasswordHash, string) error { return fake.err }

type tokenServiceFake struct {
	pair       ports.TokenPair
	claims     ports.TokenClaims
	parseKinds []bool
	issued     int
}

func (fake *tokenServiceFake) Issue(ports.AuthUser) (ports.TokenPair, error) {
	fake.issued++
	return fake.pair, nil
}

func (fake *tokenServiceFake) Parse(_ string, refresh bool) (ports.TokenClaims, error) {
	fake.parseKinds = append(fake.parseKinds, refresh)
	return fake.claims, nil
}

type blacklistFake struct {
	revoked    bool
	revokeIDs  []string
	revokeTTLs []time.Duration
}

func (fake *blacklistFake) Revoke(_ context.Context, tokenID string, ttl time.Duration) error {
	fake.revokeIDs = append(fake.revokeIDs, tokenID)
	fake.revokeTTLs = append(fake.revokeTTLs, ttl)
	return nil
}

func (fake blacklistFake) IsRevoked(context.Context, string) (bool, error) { return fake.revoked, nil }

type checkerFake struct{}

func (checkerFake) IsEmailTaken(context.Context, values.Email) (bool, error) { return false, nil }

func newAuthService(repository *authRepositoryFake, verifier verifierFake, tokens *tokenServiceFake, blacklist *blacklistFake) applicationUser.AuthService {
	return applicationUser.NewAuthService(repository, verifier, tokens, blacklist, checkerFake{})
}

func testAuthUser() ports.AuthUser {
	email, _ := values.NewEmail("john@example.com")
	hash, _ := values.NewPasswordHash("hash")
	return ports.AuthUser{ID: values.NewUserID("user-1"), Email: email, PasswordHash: hash, CreatedAt: time.Unix(0, 0).UTC()}
}

func TestAuthServiceLogin(t *testing.T) {
	tokens := &tokenServiceFake{pair: ports.TokenPair{AccessToken: "access", RefreshToken: "refresh"}}
	service := newAuthService(&authRepositoryFake{user: testAuthUser()}, verifierFake{}, tokens, &blacklistFake{})

	pair, err := service.Login(context.Background(), applicationUser.LoginRequest{Email: "john@example.com", Password: "StrongPass123"})
	require.NoError(t, err)
	assert.Equal(t, "access", pair.AccessToken)
	assert.Equal(t, 1, tokens.issued)
}

func TestAuthServiceRejectsInvalidCredentials(t *testing.T) {
	service := newAuthService(&authRepositoryFake{findErr: errors.New("not found")}, verifierFake{}, &tokenServiceFake{}, &blacklistFake{})

	_, err := service.Login(context.Background(), applicationUser.LoginRequest{Email: "unknown@example.com", Password: "bad"})
	require.ErrorIs(t, err, domainErrors.ErrInvalidCredentials)
}

func TestAuthServiceRefreshChecksBlacklistAndUser(t *testing.T) {
	tokens := &tokenServiceFake{pair: ports.TokenPair{AccessToken: "new-access", RefreshToken: "new-refresh"}, claims: ports.TokenClaims{UserID: values.NewUserID("user-1"), TokenID: "refresh-id", ExpiresAt: time.Now().Add(time.Hour), Refresh: true}}
	service := newAuthService(&authRepositoryFake{user: testAuthUser()}, verifierFake{}, tokens, &blacklistFake{})

	pair, err := service.Refresh(context.Background(), "old-refresh")
	require.NoError(t, err)
	assert.Equal(t, "new-refresh", pair.RefreshToken)
	require.Len(t, tokens.parseKinds, 1)
	assert.True(t, tokens.parseKinds[0])
}

func TestAuthServiceLogoutRevokesBothTokens(t *testing.T) {
	tokens := &tokenServiceFake{claims: ports.TokenClaims{TokenID: "token-id", ExpiresAt: time.Now().Add(time.Hour)}}
	blacklist := &blacklistFake{}
	service := newAuthService(&authRepositoryFake{}, verifierFake{}, tokens, blacklist)

	require.NoError(t, service.Logout(context.Background(), "access", "refresh"))
	assert.Len(t, blacklist.revokeIDs, 2)
	require.Len(t, tokens.parseKinds, 2)
	assert.False(t, tokens.parseKinds[0])
	assert.True(t, tokens.parseKinds[1])
}

func TestAuthServiceUpdatesProfile(t *testing.T) {
	repository := &authRepositoryFake{user: testAuthUser()}
	service := newAuthService(repository, verifierFake{}, &tokenServiceFake{}, &blacklistFake{})

	profile, err := service.UpdateProfile(context.Background(), values.NewUserID("user-1"), applicationUser.UpdateProfileRequest{Email: "JANE@EXAMPLE.COM"})
	require.NoError(t, err)
	assert.Equal(t, "jane@example.com", profile.Email.String())
	assert.Equal(t, "user-1", repository.updatedID.String())
}
