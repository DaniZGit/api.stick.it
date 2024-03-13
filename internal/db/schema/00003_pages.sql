-- +goose Up
CREATE TABLE pages (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	sort_order INTEGER NOT NULL,
	album_id UUID NOT NULL,
	file_id UUID
);

-- +goose Down
DROP TABLE IF EXISTS pages;