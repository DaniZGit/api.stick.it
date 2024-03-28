-- name: CreatePage :one
INSERT INTO pages(id, sort_order, album_id, file_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetPage :many
SELECT 
  p.*, -- page
  pf.id AS page_file_id, pf.name AS page_file_name, pf.path AS page_file_path, -- page file
  s.id AS sticker_id, s.created_at AS sticker_created_at, s.title AS sticker_title, s.type AS sticker_type, 
  s.top AS sticker_top, s.left AS sticker_left, s.width AS sticker_width, s.height AS sticker_height, 
  s.numerator AS sticker_numerator, s.denominator AS sticker_denominator,
  s.rotation AS sticker_rotation, s.rarity_id AS sticker_rarity_id, -- sticker
  sf.id AS sticker_file_id, sf.name AS sticker_file_name, sf.path AS sticker_file_path -- sticker file
FROM pages AS p
LEFT JOIN files AS pf ON p.file_id = pf.id
LEFT JOIN stickers AS s ON p.id = s.page_id
LEFT JOIN files AS sf ON s.file_id = sf.id
WHERE p.id = $1;

-- name: GetPages :many
SELECT
  p.*,
  pf.id AS page_file_id, pf.name AS page_file_name, pf.path AS page_file_path -- album file
FROM pages p
LEFT JOIN files pf ON p.file_id = pf.id
WHERE p.album_id = $1
ORDER BY sort_order ASC;