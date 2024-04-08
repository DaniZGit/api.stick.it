-- name: CreatePack :one
INSERT INTO packs(id, title, price, album_id, file_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdatePack :one
UPDATE packs SET
  title = $1,
  price = $2,
  file_id = $3
WHERE id = $4
RETURNING *;

-- name: DeletePack :one
DELETE FROM packs
WHERE id = $1
RETURNING *;

-- name: GetAlbumPacks :many
SELECT 
  p.*, -- pack
  pf.id AS pack_file_id, pf.name AS pack_file_name, pf.path AS pack_file_path, -- pack file
  pr.id AS pack_rarities_id, pr.pack_id AS pack_rarities_pack_id, pr.rarity_id AS pack_rarities_rarity_id, pr.drop_chance AS pack_rarities_drop_chance -- pack rarities
FROM packs AS p
LEFT JOIN pack_rarities AS pr ON p.id = pr.pack_id
LEFT JOIN files AS pf ON p.file_id = pf.id
WHERE p.album_id = $1;

-- name: GetPackRarities :many
SELECT
 pr.* -- pack rarities
FROM pack_rarities AS pr
WHERE pr.pack_id = $1;

-- name: CreatePackRarity :one
INSERT INTO pack_rarities(id, pack_id, rarity_id, drop_chance)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdatePackRarity :one
UPDATE pack_rarities SET
  drop_chance = $1
WHERE id = $2
RETURNING *;

-- name: DeletePackRarity :one
DELETE FROM pack_rarities
WHERE id = $1
RETURNING *;