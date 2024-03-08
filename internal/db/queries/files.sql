-- name: CreateFile :one
INSERT INTO files (id, name, path)
VALUES ($1, $2, $3)
RETURNING *;