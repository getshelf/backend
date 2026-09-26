package collection

import "time"

type Collection struct {
	ID        string
	Title     string
	Icon      *string
	ParentID  *string
	OwnerID   string
	SortOrder int
	Collections []Collection
	CreatedAt time.Time
	UpdatedAt time.Time
}
