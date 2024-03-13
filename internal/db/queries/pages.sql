-- name: CreatePage :one
INSERT INTO pages(id, sort_order, album_id, file_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetPage :many
SELECT 
  p.*, -- page
  pf.id AS page_file_id, pf.name AS page_file_name, pf.path AS page_file_path, -- page file
  s.id AS sticker_id, s.created_at AS sticker_created_at, s.title AS sticker_title, -- sticker
  sf.id AS sticker_file_id, sf.name AS sticker_file_name, sf.path AS sticker_file_path -- sticker file
FROM pages p
LEFT JOIN files pf ON p.file_id = pf.id
LEFT JOIN stickers s ON p.id = s.id
LEFT JOIN files sf ON s.file_id = sf.id
WHERE pages.id = $1;

-- name: GetPages :many
SELECT
  p.*,
  pf.id AS page_file_id, pf.name AS page_file_name, pf.path AS page_file_path -- album file
FROM pages p
LEFT JOIN files pf ON p.file_id = pf.id
WHERE p.album_id = $1
ORDER BY sort_order ASC;