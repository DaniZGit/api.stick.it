-- +goose Up
CREATE TABLE stickers (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	title VARCHAR(255) NOT NULL,
  "left" DECIMAL NOT NULL,
  "top" DECIMAL NOT NULL,
  "rotation" DECIMAL NOT NULL,
  "type" VARCHAR(255) NOT NULL,
  "width" DECIMAL NOT NULL,
  "height" DECIMAL NOT NULL,
  "numerator" INT NOT NULL,
  "denominator" INT NOT NULL, 
  page_id UUID NOT NULL,
  rarity_id UUID,
  file_id UUID,
  sticker_id UUID
);

-- +goose Down
DROP TABLE IF EXISTS stickers;