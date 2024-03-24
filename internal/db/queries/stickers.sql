-- name: CreateSticker :one
INSERT INTO stickers(id, title, "type", "top", "left", file_id, page_id, rarity_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPageStickers :many
SELECT 
  s.*, -- stickers
  sf.id AS sticker_file_id, sf.name AS sticker_file_name, sf.path AS sticker_file_path -- sticker file
FROM stickers s
LEFT JOIN files sf ON s.file_id = sf.id
WHERE s.page_id = $1;

-- name: UpdateSticker :one
UPDATE stickers
SET title = $1,
    "type" = $2,
    "top" = $3,
    "left" = $4,
    file_id = $5,
    rarity_id = $6
WHERE id = $7
RETURNING *;