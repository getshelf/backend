package account

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	store        Store
	hashPassword func(string) (string, error)
	newID        func() (string, error)
	now          func() time.Time
}

type RegisterInput struct {
	Email        string
	Password     string
}

type RegisterOutput struct {
	ID           string
}

func NewService(store Store) *Service {
	return newService(
		store,
		bcryptHash,
		randomID,
		time.Now,
	)
}

func newService(
	store Store,
	hashPassword func(string) (string, error),
	newID func() (string, error),
	now func() time.Time,
) *Service {
	return &Service{
		store:        store,
		hashPassword: hashPassword,
		newID:        newID,
		now:          now,
	}
}

func (service *Service) Register(
	ctx context.Context,
	input RegisterInput,
) (RegisterOutput, error) {
	email, err := normalizeEmail(input.Email)
	if err != nil {
		return RegisterOutput{}, err
	}

	if err := validatePassword(input.Password, email); err != nil {
		return RegisterOutput{}, err
	}

	passwordHash, err := service.hashPassword(input.Password)
	if err != nil {
		return RegisterOutput{}, fmt.Errorf("hash password: %w", err)
	}

	id, err := service.newID()
	if err != nil {
		return RegisterOutput{}, fmt.Errorf("new id: %w", err)
	}

	err = service.store.Create(ctx, CreateAccountParams{
		ID: id,
		Email: email,
		PasswordHash: passwordHash,
		CreatedAt: service.now().UTC(),
		UpdatedAt: service.now().UTC(),
	})


	if err != nil {
		if errors.Is(err, ErrEmailTaken) {
			return RegisterOutput{}, ErrEmailTaken
		}

		return RegisterOutput{}, fmt.Errorf("create account: %w", err)
	}

	return RegisterOutput{ID: id}, nil
}

func normalizeEmail(email string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	address, err := mail.ParseAddress(email)

	if err != nil || address.Address != email {
		return "", ErrInvalidEmail
	}

	return email, nil
}

func validatePassword(password, email string) error {
	if utf8.RuneCountInString(password) < 8 {
		return ErrPasswordTooWeak
	}

	return nil
}

func bcryptHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func randomID() (string, error) {
	var value [16]byte

	if _, err := rand.Read(value[:]); err != nil {
		return "", err
	}

	return hex.EncodeToString(value[:]), nil
}
