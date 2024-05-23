-- name: CreateAvatar :one
INSERT INTO avatars(id, title, file_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateAvatar :one
UPDATE avatars SET
  title = $2,
  file_id = $3
WHERE id = $1
RETURNING *;

-- name: DeleteAvatar :one
DELETE FROM avatars
WHERE id = $1
RETURNING *;

-- name: GetAvatars :many
SELECT a.*, COUNT(*) OVER() as "total_rows",
  af.id AS avatar_file_id, af.name AS avatar_file_name, af.path AS avatar_file_path -- avatar file
FROM avatars a
INNER JOIN files af on a.file_id = af.id
LIMIT $1 OFFSET $2;