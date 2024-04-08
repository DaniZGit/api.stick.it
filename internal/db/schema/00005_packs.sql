-- +goose Up
CREATE TABLE packs (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	title VARCHAR(255) NOT NULL,
	price INTEGER NOT NULL,
  album_id UUID NOT NULL,
  rarity_id UUID,
  file_id UUID
);

-- +goose Down
DROP TABLE IF EXISTS packs;