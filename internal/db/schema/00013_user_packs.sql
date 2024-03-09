-- +goose Up
CREATE TABLE user_packs (
	id UUID PRIMARY KEY,
	user_id UUID NOT NULL,
	pack_id UUID NOT NULL,
	amount INTEGER NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS user_packs;