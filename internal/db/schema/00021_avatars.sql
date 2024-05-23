-- +goose Up
CREATE TABLE avatars (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	title VARCHAR(255) NOT NULL,
  file_id UUID
);

-- +goose Down
DROP TABLE IF EXISTS avatars;