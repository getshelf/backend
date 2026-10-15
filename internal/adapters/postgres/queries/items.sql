-- name: CreateSavedURL :one
WITH saved_item AS (
    INSERT INTO items (
        id,
        owner_id,
        kind,
        source_url,
        collection_id,
        created_at,
        updated_at
    ) VALUES (
        sqlc.arg(id),
        sqlc.arg(owner_id),
        'link',
        sqlc.arg(source_url),
        sqlc.narg(collection_id),
        sqlc.arg(created_at),
        sqlc.arg(updated_at)
    )
    RETURNING *
)
SELECT
    i.*,
    CASE
        WHEN c.id IS NULL THEN NULL::jsonb
        ELSE to_jsonb(c)
    END AS collection
FROM saved_item AS i
LEFT JOIN collections AS c
    ON c.id = i.collection_id
   AND c.owner_id = i.owner_id;

-- name: CreatePendingItemExtraction :exec
INSERT INTO item_extractions (
    item_id,
    status,
    created_at,
    updated_at
) VALUES (
    sqlc.arg(item_id),
    'pending',
    sqlc.arg(created_at),
    sqlc.arg(updated_at)
);

-- name: GetItem :one
SELECT *
FROM items
WHERE items.id = sqlc.arg(id)
    AND items.owner_id = sqlc.arg(owner_id);

-- name: ListItems :many
SELECT *
FROM items
ORDER BY items.created_at DESC, items.id DESC
LIMIT sqlc.arg(page_size)
OFFSET sqlc.arg(page_offset);
