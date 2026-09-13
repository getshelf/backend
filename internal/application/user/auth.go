package user

import (
	"context"
	"time"

	domainErrors "github.com/getshelf/backend/internal/domain/user/errors"
	"github.com/getshelf/backend/internal/domain/user/ports"
	"github.com/getshelf/backend/internal/domain/user/values"
)

type AuthService struct {
	repository ports.AuthRepository
	hasher     ports.PasswordVerifier
	tokens     ports.TokenService
	blacklist  ports.TokenBlacklist
	checker    ports.UniquenessChecker
}

func NewAuthService(repository ports.AuthRepository, hasher ports.PasswordVerifier, tokens ports.TokenService, blacklist ports.TokenBlacklist, checker ports.UniquenessChecker) AuthService {
	return AuthService{repository: repository, hasher: hasher, tokens: tokens, blacklist: blacklist, checker: checker}
}

type LoginRequest struct {
	Email    string
	Password string
}

func (service AuthService) Login(ctx context.Context, request LoginRequest) (ports.TokenPair, error) {
	user, err := service.repository.FindByLogin(ctx, request.Email)
	if err != nil {
		return ports.TokenPair{}, domainErrors.ErrInvalidCredentials
	}
	if err := service.hasher.Compare(user.PasswordHash, request.Password); err != nil {
		return ports.TokenPair{}, domainErrors.ErrInvalidCredentials
	}
	return service.tokens.Issue(user)
}

func (service AuthService) Refresh(ctx context.Context, refreshToken string) (ports.TokenPair, error) {
	claims, err := service.tokens.Parse(refreshToken, true)
	if err != nil {
		return ports.TokenPair{}, domainErrors.ErrUnauthorized
	}
	revoked, err := service.blacklist.IsRevoked(ctx, claims.TokenID)
	if err != nil || revoked {
		return ports.TokenPair{}, domainErrors.ErrUnauthorized
	}
	user, err := service.repository.FindByID(ctx, claims.UserID)
	if err != nil {
		return ports.TokenPair{}, domainErrors.ErrUnauthorized
	}
	return service.tokens.Issue(user)
}

func (service AuthService) Logout(ctx context.Context, accessToken string, refreshToken string) error {
	for index, rawToken := range []string{accessToken, refreshToken} {
		if rawToken == "" {
			continue
		}
		claims, err := service.tokens.Parse(rawToken, index == 1)
		if err != nil {
			continue
		}
		if err := service.blacklist.Revoke(ctx, claims.TokenID, time.Until(claims.ExpiresAt)); err != nil {
			return err
		}
	}
	return nil
}

type Profile struct {
	ID        values.UserID
	Email     values.Email
	CreatedAt time.Time
}

func (service AuthService) Profile(ctx context.Context, id values.UserID) (Profile, error) {
	user, err := service.repository.FindByID(ctx, id)
	if err != nil {
		return Profile{}, domainErrors.ErrUserNotFound
	}
	return profileFromUser(user), nil
}

type UpdateProfileRequest struct {
	Email string
}

func (service AuthService) UpdateProfile(ctx context.Context, id values.UserID, request UpdateProfileRequest) (Profile, error) {
	current, err := service.repository.FindByID(ctx, id)
	if err != nil {
		return Profile{}, domainErrors.ErrUserNotFound
	}
	email, err := values.NewEmail(request.Email)
	if err != nil {
		return Profile{}, err
	}
	if email != current.Email {
		taken, err := service.checker.IsEmailTaken(ctx, email)
		if err != nil {
			return Profile{}, err
		}
		if taken {
			return Profile{}, domainErrors.ErrEmailTaken
		}
	}
	if err := service.repository.UpdateProfile(ctx, id, email); err != nil {
		return Profile{}, err
	}
	current.Email = email
	return profileFromUser(current), nil
}

func profileFromUser(user ports.AuthUser) Profile {
	return Profile{ID: user.ID, Email: user.Email, CreatedAt: user.CreatedAt}
}
