-- +goose Up
CREATE TABLE users (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	username VARCHAR(255) UNIQUE NOT NULL,
	email VARCHAR(255) UNIQUE NOT NULL,
	password CHAR(60) NOT NULL,
	tokens BIGINT NOT NULL DEFAULT 0,
	available_free_packs INT NOT NULL DEFAULT 1,
	last_free_pack_obtain_date TIMESTAMP NOT NULL DEFAULT NOW(),
	version BIGINT NOT NULL DEFAULT 0,
	file_id UUID,
	role_id UUID
);

-- +goose Down
DROP TABLE IF EXISTS users;