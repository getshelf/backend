package item

import (
	"encoding/json"
	"time"

	"github.com/getshelf/backend/internal/modules/collection"
)

type Item struct {
	ID          string
	OwnerID     string
	Kind        string
	Title       *string
	SourceURL   *string
	ContentText *string
	Collection  *collection.Collection
	Metadata    json.RawMessage
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
