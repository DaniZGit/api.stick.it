-- name: GetAlbums :many
SELECT albums.*, sqlc.embed(albumfiles)
FROM albums
LEFT JOIN albumfiles ON albums.file_id = albumfiles.id
LIMIT $1;

-- name: GetAlbum :one
SELECT albums.*, sqlc.embed(albumfiles)
FROM albums
LEFT JOIN albumfiles ON albums.file_id = albumfiles.id
WHERE albums.id = $1
LIMIT 1;

-- name: CreateAlbum :one
INSERT INTO albums(id, title, date_from, date_to, featured, file_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateAlbum :one
UPDATE albums SET
  title = $2,
  date_from = $3,
  date_to = $4,
  featured = $5,
  file_id = $6
WHERE id = $1
RETURNING *;

-- name: DeleteAlbum :one
DELETE FROM albums
WHERE id = $1
RETURNING *;