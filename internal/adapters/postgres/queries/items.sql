-- name: CreateSavedURL :one
INSERT INTO items (
    id,
    owner_id,
    kind,
    source_url,
    created_at,
    updated_at
) VALUES (
    sqlc.arg(id),
    sqlc.arg(owner_id),
    'link',
    sqlc.arg(source_url),
    sqlc.arg(created_at),
    sqlc.arg(updated_at)
)
RETURNING *;

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
WHERE id = sqlc.arg(id)
    AND owner_id = sqlc.arg(owner_id);

-- name: ListItems :many
SELECT *
FROM items
WHERE owner_id = sqlc.arg(owner_id)
ORDER BY created_at DESC, id DESC
LIMIT sqlc.arg(page_size)
OFFSET sqlc.arg(page_offset);
