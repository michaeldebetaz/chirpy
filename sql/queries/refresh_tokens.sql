-- name: CreateRefreshToken :one
INSERT INTO
    refresh_tokens (token, user_id, expires_at)
VALUES
    ($1, $2, $3)
RETURNING
    *;

-- name: RevokeRefreshToken :exec
UPDATE
    refresh_tokens
SET
    revoked_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE
    token = $1;
