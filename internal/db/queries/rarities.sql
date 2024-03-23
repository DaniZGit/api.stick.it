-- name: CreateRarity :one
INSERT INTO rarities(id, title)
VALUES ($1, $2)
RETURNING *;

-- name: GetRarity :one
SELECT *
FROM rarities
WHERE id = $1;

-- name: GetRarities :many
SELECT *, COUNT(*) OVER() as "total_rows"
FROM rarities
LIMIT $1 OFFSET $2;