package item

import (
	"context"
	"time"
)

type SaveURLParams struct {
	ID        string
	OwnerID   string
	URL       string
	CreatedAt time.Time
}

type Store interface {
	SaveURL(ctx context.Context, params SaveURLParams) (Item, error)
}
