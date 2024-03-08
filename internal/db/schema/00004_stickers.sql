-- +goose Up
CREATE TABLE stickers (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	title VARCHAR(255) NOT NULL,
  "left" DECIMAL NOT NULL,
  "top" DECIMAL NOT NULL,
  "type" VARCHAR(255) NOT NULL,
  album_page_id UUID NOT NULL,
  rarity_id UUID NOT NULL,
  file_id UUID
);

-- +goose Down
DROP TABLE stickers;