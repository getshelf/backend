package unit

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainErrors "github.com/getshelf/backend/internal/domain/user/errors"
	"github.com/getshelf/backend/internal/domain/user/events"
	"github.com/getshelf/backend/internal/domain/user/ports"
	userService "github.com/getshelf/backend/internal/domain/user/services"
	"github.com/getshelf/backend/internal/domain/user/values"
)

type uniquenessFake struct {
	emailTaken bool
}

func (fake uniquenessFake) IsEmailTaken(context.Context, values.Email) (bool, error) {
	return fake.emailTaken, nil
}

type repositoryFake struct {
	saved ports.User
}

func (fake *repositoryFake) Save(_ context.Context, user ports.User) error {
	fake.saved = user
	return nil
}

type hasherFake struct{}

func (hasherFake) Hash(_ context.Context, password values.PlainPassword) (values.PasswordHash, error) {
	return values.NewPasswordHash("hash:" + password.String())
}

type idGeneratorFake struct{}

func (idGeneratorFake) NewUserID(context.Context) (values.UserID, error) {
	return values.NewUserID("user-1"), nil
}

type clockFake struct{}

func (clockFake) Now() time.Time {
	return time.Date(2026, 9, 13, 12, 0, 0, 0, time.UTC)
}

type transactionFake struct{}

func (transactionFake) Do(ctx context.Context, callback func(context.Context) error) error {
	return callback(ctx)
}

type publisherFake struct {
	published events.UserRegistered
}

func (fake *publisherFake) Publish(_ context.Context, event events.UserRegistered) error {
	fake.published = event
	return nil
}

func newService(checker ports.UniquenessChecker, repository *repositoryFake, publisher ports.EventPublisher) userService.RegistrationService {
	return userService.NewRegistrationService(
		checker,
		repository,
		hasherFake{},
		idGeneratorFake{},
		clockFake{},
		transactionFake{},
		publisher,
	)
}

func TestRegistrationServiceRegister(t *testing.T) {
	repository := &repositoryFake{}
	publisher := &publisherFake{}
	service := newService(uniquenessFake{}, repository, publisher)

	id, err := service.Register(context.Background(), "JOHN@EXAMPLE.COM", "StrongPass123")
	require.NoError(t, err)
	assert.Equal(t, "user-1", id.String())
	require.NotNil(t, repository.saved)
	assert.Equal(t, "john@example.com", repository.saved.Email().String())
	assert.Equal(t, "hash:StrongPass123", repository.saved.PasswordHash().String())
	assert.Equal(t, "user-1", publisher.published.UserID.String())
}

func TestRegistrationServiceRejectsWeakPassword(t *testing.T) {
	repository := &repositoryFake{}
	service := newService(uniquenessFake{}, repository, nil)

	_, err := service.Register(context.Background(), "john@example.com", "password")
	require.ErrorIs(t, err, domainErrors.ErrPasswordTooWeak)
	assert.Nil(t, repository.saved)
}

func TestRegistrationServiceRejectsTakenEmail(t *testing.T) {
	repository := &repositoryFake{}
	service := newService(uniquenessFake{emailTaken: true}, repository, nil)

	_, err := service.Register(context.Background(), "john@example.com", "StrongPass123")
	require.ErrorIs(t, err, domainErrors.ErrEmailTaken)
	assert.Nil(t, repository.saved)
}
