-- +goose Up
CREATE TABLE opened_packs (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	title VARCHAR(255) NOT NULL,
	pack_id UUID NOT NULL,
	user_id UUID NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS opened_packs;