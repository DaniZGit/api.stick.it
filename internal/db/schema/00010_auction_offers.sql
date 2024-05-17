-- +goose Up
CREATE TABLE auction_offers (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  starting_bid INTEGER NOT NULL,
  duration INTEGER NOT NULL DEFAULT 28800000, -- 8 hours
  user_sticker_id UUID NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS auction_offers;