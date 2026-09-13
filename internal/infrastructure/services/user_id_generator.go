package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/getshelf/backend/internal/domain/user/values"
)

type RandomUserIDGenerator struct{}

func (RandomUserIDGenerator) NewUserID(_ context.Context) (values.UserID, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return values.NewUserID(hex.EncodeToString(bytes)), nil
}
