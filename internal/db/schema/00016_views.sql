-- +goose Up
CREATE VIEW albumfiles AS (
  SELECT files.*
  FROM albums
  LEFT JOIN files ON albums.file_id = files.id
);

-- +goose Down
DROP VIEW IF EXISTS albumfiles;
