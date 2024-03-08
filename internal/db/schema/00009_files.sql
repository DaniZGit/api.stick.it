-- +goose Up
CREATE TABLE files (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  name VARCHAR(255) NOT NULL,
	path TEXT NOT NULL
);

-- +goose Down
DROP TABLE files;