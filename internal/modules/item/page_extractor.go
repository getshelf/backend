package item

import "context"

type PageMetadata struct {
	Title        string
	Description  string
	ImageURL     string
	SiteName     string
	CanonicalURL string
	FinalURL     string
}

type PageExtractor interface {
	Extract(ctx context.Context, rawURL string) (PageMetadata, error)
}
