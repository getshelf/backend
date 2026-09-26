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

func (service *Service) GetCollection(
	ctx context.Context,
	collectionID string,
	ownerID string,
) (Collection, error) {
	collection, err := service.store.Get(ctx, collectionID, ownerID)
	if err != nil {
		return Collection{}, err
	}
	return collection, nil
}

func (service *Service) ListCollections(
	ctx context.Context,
	ownerID string,
) ([]Collection, error) {
	collections, err := service.store.List(ctx, &ownerID)
	if err != nil {
		return nil, err
	}

	return collections, nil
}

func (service *Service) CreateCollection(
	ctx context.Context,
	input CreateCollectionInput,
) (Collection, error) {
	id, err := service.newID()

	fmt.Println("NEW ID", id)

	if err != nil {
		return Collection{}, err
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

	return collection, nil
}

func (service *Service) UpdateCollection(
	ctx context.Context,
	collectionID string,
	ownerID string,
	params UpdateCollectionParams,
) (Collection, error) {
	existingCollection, err := service.store.Get(ctx, collectionID, ownerID)

	if err != nil {
		return Collection{}, err
	}

	if params.Title == "" {
		params.Title = existingCollection.Title
	}
	if params.Icon == nil {
		params.Icon = existingCollection.Icon
	}
	if !params.ParentID.Set {
		params.ParentID.Set = true
		params.ParentID.Value = existingCollection.ParentID
	}

	collection, err := service.store.Update(
		ctx,
		collectionID,
		ownerID,
		UpdateCollectionParams{
			Title:     params.Title,
			Icon:      params.Icon,
			ParentID:  params.ParentID,
			UpdatedAt: time.Now(),
		},
	)
	if err != nil {
		return Collection{}, err
	}

	return collection, nil
}

func randomID() (string, error) {
	var value [16]byte

	if _, err := rand.Read(value[:]); err != nil {
		return "", err
	}

	return hex.EncodeToString(value[:]), nil
}
