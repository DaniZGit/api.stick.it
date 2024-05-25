-- name: GetAlbums :many
SELECT 
  a.*, -- album
  af.id AS album_file_id, af.name AS album_file_name, af.path AS album_file_path -- album file
FROM albums AS a
LEFT JOIN files AS af ON a.file_id = af.id
LIMIT $1;

-- name: GetAlbum :many
SELECT 
  a.*, -- album
  af.id AS album_file_id, af.name AS album_file_name, af.path AS album_file_path, -- album file
  p.id AS page_id, p.created_at AS page_created_at, p.sort_order AS page_sort_order, -- page
  pf.id AS page_file_id, pf.name AS page_file_name, pf.path AS page_file_path -- page file
FROM albums AS a
LEFT JOIN files AS af ON a.file_id = af.id
LEFT JOIN pages AS p ON a.id = p.album_id
LEFT JOIN files AS pf ON p.file_id = pf.id
WHERE a.id = $1
ORDER BY p.sort_order ASC;

-- name: CreateAlbum :one
INSERT INTO albums(id, title, date_from, date_to, featured, page_numerator, page_denominator, file_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateAlbum :one
UPDATE albums SET
  title = $2,
  date_from = $3,
  date_to = $4,
  featured = $5,
  page_numerator = $6,
  page_denominator = $7,
  file_id = $8
WHERE id = $1
RETURNING *;

-- name: DeleteAlbum :one
DELETE FROM albums
WHERE id = $1
RETURNING *;

-- name: GetFeaturedAlbums :many
SELECT a.*,
  af.id AS album_file_id, af.name AS album_file_name, af.path AS album_file_path -- album file
FROM albums a
LEFT JOIN files AS af ON a.file_id = af.id
WHERE featured = true;