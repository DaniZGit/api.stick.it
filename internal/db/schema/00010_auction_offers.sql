-- +goose Up
CREATE TABLE auction_offers (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  starting_price MONEY NOT NULL,
  user_id UUID NOT NULL,
  sticker_id UUID NOT NULL
);

-- +goose Down
DROP TABLE auction_offers;