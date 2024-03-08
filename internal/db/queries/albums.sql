-- name: GetAlbums :many
SELECT albums.*, sqlc.embed(albumfiles)
FROM albums
LEFT JOIN albumfiles ON albums.file_id = albumfiles.id
LIMIT $1;

-- name: GetAlbum :one
SELECT albums.*, sqlc.embed(albumfiles)
FROM albums
LEFT JOIN albumfiles ON albums.file_id = albumfiles.id
WHERE LOWER(albums.title) = LOWER($1)
LIMIT 1;

-- name: CreateAlbum :one
INSERT INTO albums(id, title, date_from, date_to, file_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteAlbum :one
DELETE FROM albums
WHERE LOWER(title) = LOWER($1)
RETURNING *;