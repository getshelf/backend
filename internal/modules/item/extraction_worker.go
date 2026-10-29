package item

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

var (
	ErrNoPendingExtraction = errors.New("no pending extraction")
	ErrExtractionLeaseLost = errors.New("extraction lease lost")
)

type ExtractionJob struct {
	ItemID       string
	SourceURL    string
	AttemptCount int32
}

type ExtractionStore interface {
	ClaimExtraction(ctx context.Context) (ExtractionJob, error)

	CompleteExtraction(
		ctx context.Context,
		job ExtractionJob,
		metadata PageMetadata,
	) error

	FailExtraction(
		ctx context.Context,
		job ExtractionJob,
		reason string,
	) error
}

type ExtractionWorker struct {
	store     ExtractionStore
	extractor PageExtractor
	logger    *slog.Logger
}

func NewExtractionWorker(
	store ExtractionStore,
	extractor PageExtractor,
	logger *slog.Logger,
) *ExtractionWorker {
	return &ExtractionWorker{
		store:     store,
		extractor: extractor,
		logger:    logger,
	}
}

func (worker *ExtractionWorker) Run(ctx context.Context) {
	for ctx.Err() == nil {
		processed, err := worker.processNext(ctx)

		if ctx.Err() != nil {
			return
		}

		if err != nil {
			worker.logger.Error(
				"extraction worker step failed",
				"error",
				err,
			)
		}

		if processed && err == nil {
			continue
		}

		timer := time.NewTimer(time.Second)
		select {
		case <- ctx.Done():
			timer.Stop()
			return
		case <- timer.C:
		}
	}
}

func (worker *ExtractionWorker) processNext(
	ctx context.Context,
) (bool, error) {
	claimCtx, cancelClaim := context.WithTimeout(ctx, 5 * time.Second)
	job, err := worker.store.ClaimExtraction(claimCtx)
	cancelClaim()

	if errors.Is(err, ErrNoPendingExtraction) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("claim extraction: %w", err)
	}

	// This happens after the claim statement has committed.
	// No database transaction remains open during the HTTP request.
	fetchCtx, cancelFetch := context.WithTimeout(ctx, 20 * time.Second)
	metadata, fetchErr := worker.extractor.Extract(fetchCtx, job.SourceURL)
	cancelFetch()

	if ctx.Err() != nil {
		// Leave the claim for lease recovery after shutdown.
		return true, ctx.Err()
	}

	// SiteName and CanonicalURL have URL-derived fallbacks, so they
	// alone do not demonstrate that useful page metadata was found.
	if fetchErr == nil &&
		metadata.Title == "" &&
		metadata.Description == "" &&
		metadata.ImageURL == "" {
		fetchErr = errors.New("page contains no usable preview metadata")
	}

	// Use a fresh timeout so an expired fetch context doesn't prevent
	// recording the failure.
	saveCtx, cancelSave := context.WithTimeout(ctx, 5*time.Second)
	defer cancelSave()

	if fetchErr != nil {
		err := worker.store.FailExtraction(saveCtx, job, fetchErr.Error())
		if errors.Is(err, ErrExtractionLeaseLost) {
			return true, nil
		}
		if err != nil {
			return true, fmt.Errorf(
				"record extraction failure for %s: %w",
				job.ItemID,
				err,
			)
		}

		worker.logger.Warn(
			"page metadata extraction failed",
			"item_id", job.ItemID,
			"attempt", job.AttemptCount,
			"error", fetchErr,
		)
		return true, nil
	}

	err = worker.store.CompleteExtraction(saveCtx, job, metadata)

	if errors.Is(err, ErrExtractionLeaseLost) {
		return true, nil
	}

	if err != nil {
		return true, fmt.Errorf(
			"save extraction result for %s: %w",
			job.ItemID,
			err,
		)
	}

	worker.logger.Info(
		"page metadata extracted",
		"item_id", job.ItemID,
		"attempt", job.AttemptCount,
	)

	return true, nil
}
