-- name: FailExhaustedExtractions :exec
UPDATE item_extractions
SET
    status = 'failed',
    last_error = 'worker lease expired after final attempt',
    updated_at = now()
WHERE status = 'processing'
    AND attempt_count >= 3
    AND COALESCE(started_at, updated_at) <= now() - interval '2 minutes';

-- name: ClaimExtraction :one
WITH candidate AS (
    SELECT e.item_id
    FROM item_extractions AS e
    WHERE e.attempt_count < 3
        AND (
            (
                e.status = 'pending'
                AND (
                    e.attempt_count = 0
                    OR e.updated_at <= now() - interval '30 seconds'
                )
            )
            OR (
                e.status = 'processing'
                AND COALESCE(e.started_at, e.updated_at)
                    <= now() - interval '2 minutes'
            )
        )
    ORDER BY e.created_at, e.item_id
    LIMIT 1
    FOR UPDATE OF e SKIP LOCKED
),
claimed AS (
    UPDATE item_extractions as e
    SET
        status = 'processing',
        attempt_count = e.attempt_count + 1,
        extractor = 'html-metadata',
        extractor_version = '1',
        started_at = now(),
        updated_at = now()
    FROM candidate as c
    WHERE e.item_id = c.item_id
    RETURNING e.item_id, e.attempt_count
)
SELECT
    c.item_id,
    c.attempt_count,
    i.source_url
FROM claimed AS c
JOIN items AS i ON i.id = c.item_id;

-- name: CompleteExtraction :execrows
UPDATE item_extractions
SET
    status = 'ready',
    last_error = NULL,
    extracted_at = now(),
    updated_at = now()
WHERE item_id = sqlc.arg(item_id)
    AND status = 'processing'
    AND attempt_count = sqlc.arg(attempt_count)::integer
    AND started_at > now() - interval '2 minutes';

-- name: SaveExtractedMetadata :execrows
UPDATE items
SET
    title = COALESCE(
        title,
        NULLIF(sqlc.arg(title)::text, '')
    ),
    description = COALESCE(
        description,
        NULLIF(sqlc.arg(description)::text, '')
    ),
    canonical_url = COALESCE(
        NULLIF(sqlc.arg(canonical_url)::text, ''),
        canonical_url
    ),
    metadata = metadata || sqlc.arg(metadata)::jsonb,
    updated_at = now()
WHERE id = sqlc.arg(item_id);

-- name: FailExtraction :execrows
UPDATE item_extractions
SET
    status = CASE
        WHEN attempt_count >= 3 THEN 'failed'
        ELSE 'pending'
    END,
    last_error = sqlc.arg(last_error)::text,
    updated_at = now()
WHERE item_id = sqlc.arg(item_id)
    AND status = 'processing'
    AND attempt_count = sqlc.arg(attempt_count)::integer
    AND started_at > now() - interval '2 minutes';
