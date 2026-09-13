package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/getshelf/backend/internal/config"
	"github.com/getshelf/backend/internal/domain/user/ports"
	"github.com/getshelf/backend/internal/domain/user/values"
)

type JWTService struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewJWTService(cfg config.Config) JWTService {
	return JWTService{secret: []byte(cfg.JWT.Secret), accessTTL: cfg.JWT.AccessTTL, refreshTTL: cfg.JWT.RefreshTTL}
}

type claims struct {
	UserID string `json:"sub"`
	Kind   string `json:"kind"`
	jwt.RegisteredClaims
}

func (service JWTService) Issue(user ports.AuthUser) (ports.TokenPair, error) {
	access, err := service.issue(user.ID, "access", service.accessTTL)
	if err != nil {
		return ports.TokenPair{}, err
	}
	refresh, err := service.issue(user.ID, "refresh", service.refreshTTL)
	if err != nil {
		return ports.TokenPair{}, err
	}
	return ports.TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}

func (service JWTService) issue(userID values.UserID, kind string, ttl time.Duration) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{UserID: userID.String(), Kind: kind, RegisteredClaims: jwt.RegisteredClaims{
		ID: uuid(), IssuedAt: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
	}})
	return token.SignedString(service.secret)
}

func (service JWTService) Parse(raw string, refresh bool) (ports.TokenClaims, error) {
	parsed, err := jwt.ParseWithClaims(raw, &claims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return service.secret, nil
	})
	if err != nil {
		return ports.TokenClaims{}, err
	}
	parsedClaims, ok := parsed.Claims.(*claims)
	if !ok || parsedClaims.Kind != map[bool]string{true: "refresh", false: "access"}[refresh] || parsedClaims.ExpiresAt == nil {
		return ports.TokenClaims{}, errors.New("invalid token claims")
	}
	return ports.TokenClaims{UserID: values.NewUserID(parsedClaims.UserID), TokenID: parsedClaims.ID, ExpiresAt: parsedClaims.ExpiresAt.Time, Refresh: refresh}, nil
}

func uuid() string {
	return time.Now().UTC().Format("20060102150405.000000000")
}
