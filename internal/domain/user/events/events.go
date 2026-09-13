package events

import (
	"time"

	"github.com/getshelf/backend/internal/domain/user/values"
)

type UserRegistered struct {
	UserID    values.UserID
	Email     values.Email
	CreatedAt time.Time
}
