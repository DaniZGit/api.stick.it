-- +goose Up
CREATE UNIQUE INDEX albums_title_unique ON albums (UPPER(title));

-- +goose Down
DROP INDEX IF EXISTS albums_title_unique;