-- +goose Up
CREATE TABLE user_stickers (
	id UUID PRIMARY KEY,
	user_id UUID NOT NULL,
	sticker_id UUID NOT NULL,
	amount INTEGER NOT NULL,
	sticked BOOLEAN NOT NULL DEFAULT 0
);

-- +goose Down
DROP TABLE IF EXISTS user_stickers;