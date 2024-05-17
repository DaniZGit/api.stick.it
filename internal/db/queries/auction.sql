-- name: CreateAuctionOFfer :one
INSERT INTO auction_offers(id, starting_bid, user_sticker_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAuctionOffers :many
SELECT ao.*, COALESCE(ab.bid, 0) as latest_bid,
  us.user_id AS user_sticker_user_id, us.sticker_id AS user_sticker_sticker_id, us.amount AS user_sticker_amount, us.sticked AS user_sticker_sticked, -- user_sticker
  s.id AS sticker_id, s.created_at AS sticker_created_at, s.title AS sticker_title, s.type AS sticker_type, 
  s.top AS sticker_top, s.left AS sticker_left, s.width AS sticker_width, s.height AS sticker_height, 
  s.numerator AS sticker_numerator, s.denominator AS sticker_denominator,
  s.rotation AS sticker_rotation, s.page_id AS sticker_page_id, s.sticker_id AS sticker_sticker_id, s.rarity_id AS sticker_rarity_id, -- sticker
  r.id AS sticker_rarity_id, r.title AS sticker_rarity_title, -- sticker rarity
  sf.id AS sticker_file_id, sf.name AS sticker_file_name, sf.path AS sticker_file_path, -- sticker file
  a.id AS album_id, a.title AS album_title -- album
FROM auction_offers ao
INNER JOIN user_stickers us ON ao.user_sticker_id = us.id
INNER JOIN stickers s ON us.sticker_id = s.id
LEFT JOIN rarities r ON s.rarity_id = r.id
LEFT JOIN files sf ON s.file_id = sf.id
INNER JOIN pages p ON s.page_id = p.id
INNER JOIN albums a ON a.id = p.album_id
LEFT join LATERAL ( -- gets last auction bid
	select * 
	from auction_bids 
	where auction_offer_id = ao.id
	order by bid desc
	limit 1
) as ab on ab.auction_offer_id = ao.id;

-- name: GetAuctionOffer :one
SELECT *
FROM auction_offers
WHERE id = $1
LIMIT 1;

-- name: CreateAuctionBid :one
INSERT INTO auction_bids(id, bid, auction_offer_id, user_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetAuctionBids :many
SELECT ab.*,
  u.id AS user_user_id, u.username AS user_username, u.email AS user_email, u.tokens AS user_tokens, -- user
  uf.id AS user_file_id, uf.name AS user_file_name, uf.path AS user_file_path -- user file
FROM auction_bids ab
INNER JOIN users u ON ab.user_id = u.id
LEFT JOIN files uf ON uf.id = u.file_id
WHERE ab.auction_offer_id = $1
ORDER BY ab.bid ASC;

-- name: GetLatestAuctionBid :one
SELECT ab.*,
  u.id AS user_user_id, u.username AS user_username, u.email AS user_email, u.tokens AS user_tokens, -- user
  uf.id AS user_file_id, uf.name AS user_file_name, uf.path AS user_file_path -- user file
FROM auction_bids ab
INNER JOIN users u ON ab.user_id = u.id
LEFT JOIN files uf ON uf.id = u.file_id
WHERE ab.auction_offer_id = $1
ORDER BY ab.created_at DESC
LIMIT 1;