package item

import (
	"encoding/json"
	"time"
)

type Item struct {
	ID          string
	OwnerID     string
	Kind        string
	Title       *string
	SourceURL   *string
	ContentText *string
	Metadat     json.RawMessage
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
