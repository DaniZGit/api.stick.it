-- name: CreateUser :one
INSERT INTO users (id, username, email, password, role_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUser :one
SELECT *
FROM users
WHERE lower(email)=lower($1);

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id=$1;

-- name: GetUsers :many
SELECT *
FROM users
LIMIT $1 OFFSET $2;

-- name: IncrementUserTokens :one
UPDATE users
SET tokens = tokens + $1
WHERE id = $2
RETURNING *;

-- name: DecrementUserTokens :one
UPDATE users
SET tokens = tokens - $1
WHERE id = $2
RETURNING *;

-- name: GetUserPacks :many
SELECT up.*,
  p.created_at as pack_created_at, p.title as pack_title, p.price as pack_price, p.amount pack_amount, p.album_id as pack_album_id, -- pack
  pf.id AS pack_file_id, pf.name AS pack_file_name, pf.path AS pack_file_path -- pack file
FROM user_packs up
INNER JOIN packs p on up.pack_id = p.id
LEFT JOIN files pf on p.file_id = pf.id
WHERE up.user_id = $1 AND p.album_id = $2 AND up.amount > 0;

-- name: GetUserPack :one
SELECT up.*
FROM user_packs up
WHERE up.user_id = $1 AND up.pack_id = $2
LIMIT 1;

-- name: GetUserStickers :many
SELECT
  us.*,
  s.created_at AS sticker_created_at, s.title AS sticker_title, s.type AS sticker_type, 
  s.top AS sticker_top, s.left AS sticker_left, s.width AS sticker_width, s.height AS sticker_height, 
  s.numerator AS sticker_numerator, s.denominator AS sticker_denominator,
  s.rotation AS sticker_rotation, s.page_id AS sticker_page_id, s.sticker_id AS sticker_sticker_id, s.rarity_id AS sticker_rarity_id, -- sticker
  r.id AS sticker_rarity_id, r.title AS sticker_rarity_title, -- sticker rarity
  sf.id AS sticker_file_id, sf.name AS sticker_file_name, sf.path AS sticker_file_path -- sticker file
FROM user_stickers us
INNER JOIN stickers s ON us.sticker_id = s.id
LEFT JOIN rarities r ON s.rarity_id = r.id
LEFT JOIN files sf ON s.file_id = sf.id
INNER JOIN pages p ON s.page_id = p.id
WHERE us.user_id = $1 AND p.album_id = $2 AND (us.amount > 0 OR us.sticked = true)
ORDER BY us.id DESC;