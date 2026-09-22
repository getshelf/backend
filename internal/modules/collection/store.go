package collection

import (
	"context"
	"time"
)

type CreateCollectionParams struct {
	ID        string
	Title     string
	Icon      *string
	ParentID  *string
	OwnerID   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Store interface {
	Create(
		ctx context.Context,
		params CreateCollectionParams,
	) (Collection, error)
}
