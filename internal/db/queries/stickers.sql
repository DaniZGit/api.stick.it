-- name: CreateSticker :one
INSERT INTO stickers(id, title, "type", "top", "left", "width", "height", "numerator", "denominator", "rotation", file_id, page_id, rarity_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
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
    "width" = $5,
    "height" = $6,
    "numerator" = $7,
    "denominator" = $8,
    "rotation" = $9,
    file_id = $10,
    rarity_id = $11
WHERE id = $12
RETURNING *;

-- name: DeleteSticker :one
DELETE FROM stickers
WHERE id = $1
RETURNING *;