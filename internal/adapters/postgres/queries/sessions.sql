-- name: CreateSession :exec
INSERT INTO sessions (
    token_hash,
    account_id,
    expires_at,
    created_at
)
VALUES (
    $1,
    $2,
    $3,
    $4
);

-- name: FindSessionByTokenHash :one
SELECT
    token_hash,
    account_id,
    expires_at,
    created_at
FROM sessions
WHERE token_hash = $1;

-- name: DeleteSessionByTokenHash :exec
DELETE FROM sessions
WHERE token_hash = $1;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions
WHERE expires_at <= $1;
