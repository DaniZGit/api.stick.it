-- +goose Up
CREATE TABLE user_stickers (
	id UUID PRIMARY KEY,
	user_id UUID NOT NULL,
	sticker_id UUID NOT NULL,
	amount INTEGER NOT NULL
);

-- +goose Down
DROP TABLE user_stickers;