-- +goose Up
CREATE TABLE opened_pack_stickers (
	id UUID PRIMARY KEY,
	sticker_id UUID NOT NULL,
	opened_pack_id UUID NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS opened_pack_stickers;