-- +goose Up
CREATE TABLE roles (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	title VARCHAR(255) UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS roles;