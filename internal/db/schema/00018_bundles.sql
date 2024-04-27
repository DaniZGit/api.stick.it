-- +goose Up
CREATE TABLE bundles (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	title VARCHAR(255) NOT NULL,
	price DECIMAL NOT NULL,
	tokens INTEGER NOT NULL,
	bonus INTEGER NOT NULL,
  file_id UUID
);

-- +goose Down
DROP TABLE IF EXISTS bundles;