package services

import (
	"context"

	"github.com/getshelf/backend/internal/domain/user/entities"
	"github.com/getshelf/backend/internal/domain/user/events"
	"github.com/getshelf/backend/internal/domain/user/ports"
	"github.com/getshelf/backend/internal/domain/user/values"
)

type RegistrationService struct {
	checker   ports.UniquenessChecker
	repo      ports.UserRepository
	hasher    ports.PasswordHasher
	id        ports.IDGenerator
	clock     ports.Clock
	tx        ports.TxManager
	publisher ports.EventPublisher
}

func NewRegistrationService(checker ports.UniquenessChecker, repo ports.UserRepository, hasher ports.PasswordHasher, id ports.IDGenerator, clock ports.Clock, tx ports.TxManager, publisher ports.EventPublisher) RegistrationService {
	return RegistrationService{checker: checker, repo: repo, hasher: hasher, id: id, clock: clock, tx: tx, publisher: publisher}
}

func (service RegistrationService) Register(ctx context.Context, rawEmail string, rawPassword string) (values.UserID, error) {
	email, err := values.NewEmail(rawEmail)
	if err != nil {
		return "", err
	}
	password, err := values.NewPlainPassword(rawPassword, email)
	if err != nil {
		return "", err
	}
	hash, err := service.hasher.Hash(ctx, password)
	if err != nil {
		return "", err
	}
	id, err := service.id.NewUserID(ctx)
	if err != nil {
		return "", err
	}
	now := service.clock.Now()

	var user entities.User
	err = service.tx.Do(ctx, func(txContext context.Context) error {
		user, err = entities.NewUser(txContext, id, email, hash, service.checker, now)
		if err != nil {
			return err
		}
		return service.repo.Save(txContext, user)
	})
	if err != nil {
		return "", err
	}

	if service.publisher != nil {
		err = service.publisher.Publish(ctx, events.UserRegistered{UserID: user.ID(), Email: user.Email(), CreatedAt: user.CreatedAt()})
		if err != nil {
			return "", err
		}
	}
	return user.ID(), nil
}
