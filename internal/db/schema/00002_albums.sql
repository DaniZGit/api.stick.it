-- +goose Up
CREATE TABLE albums (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	title VARCHAR(255) UNIQUE NOT NULL,
  date_from TIMESTAMP NOT NULL,
  date_to TIMESTAMP,
  featured BOOLEAN,
  file_id UUID
);

-- +goose Down
DROP TABLE albums;