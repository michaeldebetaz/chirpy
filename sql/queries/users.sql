-- name: CreateUser :one
INSERT INTO
    users (email, password_hash)
VALUES
    ($1, $2)
RETURNING
    *;

-- name: GetUserByRefreshToken :one
SELECT
    sqlc.embed(users),
    sqlc.embed(refresh_tokens)
FROM
    users
    INNER JOIN refresh_tokens ON users.id = refresh_tokens.user_id
WHERE
    refresh_tokens.token = $1;

-- name: GetUserByEmail :one
SELECT
    *
FROM
    users
WHERE
    email = $1;

-- name: UpdateUser :one
UPDATE
    users
SET
    email = $1,
    password_hash = $2
WHERE
    id = $3
RETURNING
    *;

-- name: DeleteAllUsers :exec
DELETE FROM
    users;
