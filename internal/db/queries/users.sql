-- name: CreateUser :one
INSERT INTO users (id, username, email, password, role_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUser :one
SELECT *
FROM users
WHERE lower(email)=lower($1);

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id=$1;

-- name: GetUsers :many
SELECT *
FROM users
LIMIT $1 OFFSET $2;

-- name: IncrementUserTokens :one
UPDATE users
SET tokens = tokens + $1
WHERE id = $2
RETURNING *;

-- name: DecrementUserTokens :one
UPDATE users
SET tokens = tokens - $1
WHERE id = $2
RETURNING *;
