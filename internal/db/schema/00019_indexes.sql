-- +goose Up
CREATE UNIQUE INDEX albums_title_unique ON albums (UPPER(title));
CREATE UNIQUE INDEX roles_title_unique ON roles (UPPER(title));
CREATE UNIQUE INDEX pack_rarities_unique ON pack_rarities (pack_id, rarity_id);

-- +goose Down
DROP INDEX IF EXISTS albums_title_unique;
DROP INDEX IF EXISTS roles_title_unique;
DROP INDEX IF EXISTS pack_rarities_unique;