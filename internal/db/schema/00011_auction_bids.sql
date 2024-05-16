-- +goose Up
CREATE TABLE auction_bids (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  bid INTEGER NOT NULL,
  auction_offer_id UUID NOT NULL,
  user_id UUID NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS auction_bids;