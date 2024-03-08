-- +goose Up
CREATE TABLE rarities (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	title VARCHAR(255) NOT NULL
);

-- +goose Down
DROP TABLE rarities;