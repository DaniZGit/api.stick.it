-- name: CreateSticker :one
INSERT INTO stickers(id, title, "type", "top", "left", "width", "height", "numerator", "denominator", "rotation", file_id, page_id, rarity_id, sticker_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING *;

-- name: GetSticker :one
SELECT *
FROM stickers
WHERE id = $1
LIMIT 1;

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
    file_id = $10
WHERE id = $11
RETURNING *;

-- name: DeleteSticker :one
DELETE FROM stickers
WHERE id = $1
RETURNING *;

-- name: GetStickerRarities :many
SELECT 
  s.*, -- stickers
  r.id AS sticker_rarity_id, r.title AS sticker_rarity_title, -- sticker rarity
  sf.id AS sticker_file_id, sf.name AS sticker_file_name, sf.path AS sticker_file_path -- sticker file
FROM stickers s
LEFT JOIN rarities r ON s.rarity_id = r.id
LEFT JOIN files sf ON s.file_id = sf.id
WHERE sticker_id = $1 AND rarity_id IS NOT NULL
ORDER BY s.created_at ASC;

-- name: CreateUserSticker :one
INSERT INTO user_stickers(id, user_id, sticker_id, amount, sticked)
VALUES($1, $2, $3, $4, $5)
ON CONFLICT ON CONSTRAINT user_stickers_unique
DO UPDATE SET amount = user_stickers.amount + EXCLUDED.amount
RETURNING *;

-- name: StickUserSticker :one
UPDATE user_stickers 
SET sticked = true,
    amount = amount - 1
WHERE user_id = $1 AND sticker_id = $2 AND amount > 0 AND sticked = false
RETURNING *;

-- name: GetRandomStickers :many
SELECT s.*,
  r.id AS sticker_rarity_id, r.title AS sticker_rarity_title, -- sticker rarity
  sf.id AS sticker_file_id, sf.name AS sticker_file_name, sf.path AS sticker_file_path -- sticker file
FROM stickers s
INNER JOIN pages p ON s.page_id = p.id
LEFT JOIN rarities r ON s.rarity_id = r.id
LEFT JOIN files sf ON s.file_id = sf.id
WHERE p.album_id = $1 AND ( (s.rarity_id = $2) OR ($2 IS NULL AND s.rarity_id IS NULL) )
ORDER BY random()
LIMIT $3;