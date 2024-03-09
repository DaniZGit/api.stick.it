-- name: CreateFile :one
INSERT INTO files (id, name, path)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetFile :one
SELECT * FROM files
WHERE id = $1
LIMIT 1;