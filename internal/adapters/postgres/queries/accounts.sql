-- name: CreateAccount :exec
INSERT INTO accounts (
    id,
    email,
    password_hash,
    created_at,
    updated_at
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
);

-- name: FindAccountByEmail :one
SELECT
    id,
    email,
    password_hash,
    created_at,
    updated_at
FROM accounts
WHERE email = $1;
