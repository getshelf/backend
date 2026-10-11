package item

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/url"
	"strings"
	"time"
)

var ErrInvalidURL = errors.New("invalid URL")

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return newService(
		store,
	)
}

func newService(
	store Store,
) *Service {
	return &Service{
		store: store,
	}
}

func (service *Service) SaveURL(
	ctx context.Context,
	ownerID string,
	rawURL string,
) (Item, error) {
	rawURL = strings.TrimSpace(rawURL)

	parsed, err := url.Parse(rawURL)

	if err != nil {
		return Item{}, ErrInvalidURL
	}

	if parsed.Hostname() == "" ||
		(parsed.Scheme != "http" && parsed.Scheme != "https") {
		return Item{}, ErrInvalidURL
	}

	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return Item{}, err
	}

	return service.store.SaveURL(ctx, SaveURLParams{
		ID: hex.EncodeToString(bytes[:]),
		OwnerID: ownerID,
		URL: rawURL,
		CreatedAt: time.Now().UTC(),
	})
}
