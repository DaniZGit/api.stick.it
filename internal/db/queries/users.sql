-- name: CreateUser :one
INSERT INTO users (id, username, email, password, confirmation_token, role_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ConfirmUserMail :one
UPDATE users
SET confirmation_token = NULL
WHERE confirmation_token = $1
RETURNING *;

-- name: GetUser :one
SELECT *
FROM users
WHERE lower(email)=lower($1);

-- name: GetUserByID :one
SELECT u.*, 
  a.title as avatar_title, -- avatar
  af.id AS avatar_file_id, af.name AS avatar_file_name, af.path AS avatar_file_path -- avatar file
FROM users u
LEFT JOIN avatars a on u.avatar_id = a.id
LEFT JOIN files af on a.file_id = af.id
WHERE u.id=$1;

-- name: GetUsers :many
SELECT *
FROM users
LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET 
  description = $1,
  avatar_id = $2
WHERE id = $3
RETURNING *;

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

-- name: GetUserAlbums :many
SELECT a.*, 
  af.id AS album_file_id, af.name AS album_file_name, af.path AS album_file_path, -- album file
  coalesce(ups.user_packs_amount, 0) as "user_packs_amount", coalesce(uss.user_stickers_amount, 0) as "user_stickers_amount", COUNT(ps.id) as "stickers_amount"
FROM albums a
LEFT JOIN files AS af ON a.file_id = af.id
INNER JOIN pages ap on a.id = ap.album_id 
INNER JOIN stickers ps on ap.id = ps.page_id
LEFT JOIN LATERAL (
	SELECT p.album_id, SUM(up.amount) as user_packs_amount
	FROM packs p
  	INNER JOIN user_packs up ON p.id = up.pack_id
	WHERE p.album_id = a.id AND up.user_id = $1 and up.amount > 0
  	GROUP BY p.album_id
) as ups on ups.album_id = a.id
LEFT JOIN LATERAL (
	SELECT p.album_id, COUNT(us.id) as user_stickers_amount
	FROM user_stickers us 
  INNER JOIN stickers s ON us.sticker_id = s.id
  INNER JOIN pages p on s.page_id = p.id
	WHERE us.user_id = $1 and us.sticked = true
  	GROUP BY p.album_id
) as uss on uss.album_id = a.id
GROUP BY a.id, af.id, ups.user_packs_amount, uss.user_stickers_amount
ORDER BY a.created_at DESC;

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

-- name: GetUserSticker :one
SELECT *
FROM user_stickers
WHERE id = $1
LIMIT 1;

-- name: GetUserStickersForAlbum :many
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

-- name: GetUserAuctionStickers :many
SELECT
  us.*,
  s.created_at AS sticker_created_at, s.title AS sticker_title, s.type AS sticker_type, 
  s.top AS sticker_top, s.left AS sticker_left, s.width AS sticker_width, s.height AS sticker_height, 
  s.numerator AS sticker_numerator, s.denominator AS sticker_denominator,
  s.rotation AS sticker_rotation, s.page_id AS sticker_page_id, s.sticker_id AS sticker_sticker_id, s.rarity_id AS sticker_rarity_id, -- sticker
  r.id AS sticker_rarity_id, r.title AS sticker_rarity_title, -- sticker rarity
  sf.id AS sticker_file_id, sf.name AS sticker_file_name, sf.path AS sticker_file_path, -- sticker file
  a.id AS album_id, a.title AS album_title -- album
FROM user_stickers us
INNER JOIN stickers s ON us.sticker_id = s.id
LEFT JOIN rarities r ON s.rarity_id = r.id
LEFT JOIN files sf ON s.file_id = sf.id
INNER JOIN pages p ON s.page_id = p.id
INNER JOIN albums a ON a.id = p.album_id
WHERE us.user_id = $1 AND us.amount > 0 AND s.rarity_id IS NOT NULL
ORDER BY us.id DESC;

-- name: UpdateUsersFreePacks :exec
UPDATE users
SET available_free_packs = available_free_packs + 1,
last_free_pack_obtain_date = NOW()
WHERE available_free_packs < $1
AND DATE_PART('hour', AGE(NOW(), last_free_pack_obtain_date)) >= 12;

-- name: ClaimUserFreePack :one
UPDATE users
SET available_free_packs = available_free_packs - 1
WHERE id = $1 AND available_free_packs > 0
RETURNING *;

-- name: ResetUserFreePackDate :one
UPDATE users
SET last_free_pack_obtain_date = NOW()
WHERE id = $1
RETURNING *;

-- name: DecreaseUserStickerAmount :one
UPDATE user_stickers
SET amount = amount - 1
WHERE id = $1 AND amount > 0
RETURNING *;

-- name: GetUserFoundStickersCount :one
SELECT count(us.sticker_id) as stickers_found
FROM users u
inner join user_stickers us on u.id = us.user_id
WHERE u.id = $1
group by u.id
LIMIT 1;

-- name: GetUserCompletedAlbumsCount :many
select 1 as completed
from user_stickers us
inner join stickers s2 on s2.id = us.sticker_id
inner join pages p2 on p2.id = s2.page_id
where us.user_id = $1 and us.sticked = true
group by p2.album_id 
having count(us.sticker_id) = (
  select count(s.id) as stickers_count
  from albums a
  inner join pages p on a.id = p.album_id
  inner join stickers s on p.id = s.page_id
  where a.id = p2.album_id 
  group by a.id
  limit 1
);