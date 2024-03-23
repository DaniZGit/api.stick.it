-- name: CreateUser :one
INSERT INTO users (id, username, email, password, role_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUser :one
SELECT *
FROM users
WHERE email=$1;

-- name: GetUsers :many
SELECT *
FROM users
LIMIT $1 OFFSET $2;
