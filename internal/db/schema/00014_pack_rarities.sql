-- +goose Up
CREATE TABLE pack_rarities (
	id UUID PRIMARY KEY,
	pack_id UUID NOT NULL,
	rarity_id UUID,
	drop_chance DECIMAL NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS pack_rarities;