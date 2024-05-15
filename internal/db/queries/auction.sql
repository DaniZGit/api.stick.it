-- name: CreateAuctionOFfer :one
INSERT INTO auction_offers(id, starting_bid, user_sticker_id)
SELECT id, starting_bid, user_sticker_id FROM (VALUES (cast($1 as UUID), CAST($2 as INTEGER), cast($3 as UUID))) AS i(id, starting_bid, user_sticker_id)
WHERE EXISTS (
  SELECT FROM user_stickers us
  WHERE us.id = cast($3 as UUID)
  AND us.amount > 0
)
RETURNING *;

-- name: GetAuctionOffers :many
SELECT *
FROM auction_offers;
