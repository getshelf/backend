package session

import (
	"time"
)

const SessionCookieName = "shelf_session"

type Session struct {
	TokenHash []byte
	AccountID string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type CreatedSession struct {
	Token string
	ExpiresAt time.Time
}
