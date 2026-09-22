package collection

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

type Service struct {
	store Store
	newID func() (string, error)
}

func NewService(
	store Store,
) *Service {
	return &Service{
		store: store,
		newID: randomID,
	}
}

type CreateCollectionInput struct {
	Title     string
	Icon      *string
	ParentID  *string
	OwnerID   string
}

type CreateCollectionOutput struct {
	ID        string
	Title     string
	Icon      *string
	ParentID  *string
	OwnerID   string
	SortOrder int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (service *Service) CreateCollection(
	ctx context.Context,
	input CreateCollectionInput,
) (CreateCollectionOutput, error) {
	id, err := service.newID()

	fmt.Println("NEW ID", id)

	if err != nil {
		return CreateCollectionOutput{}, err
	}

	collection, err := service.store.Create(ctx, CreateCollectionParams{
		ID:        id,
		Title:     input.Title,
		Icon:      input.Icon,
		ParentID:  input.ParentID,
		OwnerID:   input.OwnerID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	return CreateCollectionOutput{
		ID:        collection.ID,
		Title:     collection.Title,
		Icon:      collection.Icon,
		ParentID:  collection.ParentID,
		SortOrder: collection.SortOrder,
		OwnerID:   collection.OwnerID,
		CreatedAt: collection.CreatedAt,
		UpdatedAt: collection.UpdatedAt,
	}, nil
}

func randomID() (string, error) {
	var value [16]byte

	if _, err := rand.Read(value[:]); err != nil {
		return "", err
	}

	return hex.EncodeToString(value[:]), nil
}
