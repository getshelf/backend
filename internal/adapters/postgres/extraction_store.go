package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/getshelf/backend/internal/adapters/postgres/sqlcgen"
	"github.com/getshelf/backend/internal/modules/item"
)


var _ item.ExtractionStore = (*ItemStore)(nil)

func (store *ItemStore) ClaimExtraction(
	ctx context.Context,
) (item.ExtractionJob, error) {
	if err := store.queries.FailExhaustedExtractions(ctx); err != nil {
		return item.ExtractionJob{}, fmt.Errorf(
			"fail exhausted extractions: %w",
			err,
		)
	}

	// ClaimExtraction is one atomic SQL statement.
	row, err := store.queries.ClaimExtraction(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return item.ExtractionJob{}, item.ErrNoPendingExtraction
	}
	if err != nil {
		return item.ExtractionJob{}, err
	}

	return item.ExtractionJob{
		ItemID: row.ItemID,
		SourceURL: row.SourceUrl.String,
		AttemptCount: row.AttemptCount,
	}, nil
}

func (store *ItemStore) CompleteExtraction(
	ctx context.Context,
	job item.ExtractionJob,
	metadata item.PageMetadata,
) error {
	extra, err := json.Marshal(map[string]any{
		"page": map[string]string{
			"siteName": metadata.SiteName,
			"imageUrl": metadata.ImageURL,
			"finalUrl": metadata.FinalURL,
		},
	})

	if err != nil {
		return fmt.Errorf("encode page metadata: %w", err)
	}

	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin extraction completion: %w", err)
	}
	defer tx.Rollback()

	queries := store.queries.WithTx(tx)

	// Check ownership of the current attempt first.
	// Updating the extraction row also locks it until commit.
	affected, err := queries.CompleteExtraction(
		ctx,
		sqlcgen.CompleteExtractionParams{
			ItemID: job.ItemID,
			AttemptCount: job.AttemptCount,
		},
	)

	if err != nil {
		return fmt.Errorf("complete extraction: %w", err)
	}

	if affected != 1 {
		return item.ErrExtractionLeaseLost
	}

	affected, err = queries.SaveExtractedMetadata(
		ctx,
		sqlcgen.SaveExtractedMetadataParams{
			ItemID:       job.ItemID,
			Title:        metadata.Title,
			Description:  metadata.Description,
			CanonicalUrl: metadata.CanonicalURL,
			Metadata:     json.RawMessage(extra),
		},
	)

	if err != nil {
		return fmt.Errorf("save extracted metadata: %w", err)
	}

	if affected != 1 {
		return fmt.Errorf("item %s no longer exists", job.ItemID)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit extraction completion: %w", err)
	}

	return nil
}

func (store *ItemStore) FailExtraction(
	ctx context.Context,
	job item.ExtractionJob,
	reason string,
) error {
	// Keep stored errors bounded.
	runes := []rune(reason)
	if len(runes) > 2000 {
		reason = string(runes[:2000])
	}

	affected, err := store.queries.FailExtraction(
		ctx,
		sqlcgen.FailExtractionParams{
			ItemID: job.ItemID,
			AttemptCount: job.AttemptCount,
			LastError: reason,
		},
	)

	if err != nil {
		return err
	}

	if affected != 1 {
		return item.ErrExtractionLeaseLost
	}

	return nil
}
