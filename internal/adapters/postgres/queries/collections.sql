-- name: GetCollection :one
SELECT *
FROM collections
WHERE id = $1
  AND owner_id = $2;

-- name: CreateCollection :one
INSERT INTO collections (
    id,
    title,
    icon,
    parent_id,
    owner_id,
    sort_order,
    created_at,
    updated_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    (
        SELECT COALESCE(MAX(sort_order), -1) + 1
        FROM collections
        WHERE owner_id = $5
    ),
    $6,
    $7
) RETURNING *;

-- name: UpdateCollection :one
UPDATE collections
SET
    title = $1,
    icon = $2,
    parent_id = $3,
    updated_at = $4
WHERE id = $5
  AND owner_id = $6
RETURNING *;

-- name: ListCollections :many
SELECT *
FROM collections
WHERE owner_id = $1
ORDER BY sort_order;
