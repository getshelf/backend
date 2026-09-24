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

type ParentID struct {
	Set   bool
	Value *string
}

type UpdateCollectionParams struct {
	Title     string
	Icon      *string
	ParentID  ParentID
	UpdatedAt time.Time
}

type Store interface {
	Get(
		ctx context.Context,
		collectionID string,
		ownerID string,
	) (Collection, error)
	Create(
		ctx context.Context,
		params CreateCollectionParams,
	) (Collection, error)
	Update(
		ctx context.Context,
		collectionID string,
		ownerID string,
		params UpdateCollectionParams,
	) (Collection, error)
}
