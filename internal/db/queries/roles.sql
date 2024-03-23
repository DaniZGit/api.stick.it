-- name: CreateRole :one
INSERT INTO roles(id, title)
VALUES ($1, $2)
RETURNING *;

-- name: GetRole :one
SELECT *
FROM roles
WHERE id = $1;

-- name: GetRoleByName :one
SELECT *
FROM roles
WHERE LOWER(title) = LOWER($1);

-- name: GetRoles :many
SELECT *, COUNT(*) OVER() as "total_rows"
FROM roles
LIMIT $1 OFFSET $2;

-- name: GetRolesCount :one
SELECT COUNT(id)
FROM roles;