-- name: GetBundles :many
SELECT 
  b.*, -- bundle
  bf.id AS bundle_file_id, bf.name AS bundle_file_name, bf.path AS bundle_file_path, -- bundle file
  COUNT(*) OVER() as "total_rows"
FROM bundles AS b
LEFT JOIN files bf ON b.file_id = bf.id
ORDER BY b.price ASC
LIMIT $1 OFFSET $2;

-- name: CreateBundle :one
INSERT INTO bundles(id, title, price, tokens, bonus, file_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateBundle :one
UPDATE bundles SET
  title = $1,
  price = $2,
  tokens = $3,
  bonus = $4,
  file_id = $5
WHERE id = $6
RETURNING *;

-- name: DeleteBundle :one
DELETE FROM bundles
WHERE id = $1
RETURNING *;