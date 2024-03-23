-- +goose Up
CREATE UNIQUE INDEX albums_title_unique ON albums (UPPER(title));
CREATE UNIQUE INDEX roles_title_unique ON roles (UPPER(title));

-- +goose Down
DROP INDEX IF EXISTS albums_title_unique;
DROP INDEX IF EXISTS roles_title_unique;