package account

import "time"

type Account struct {
	ID        string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
