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
