-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, name, password)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateRefreshToken :exec
UPDATE users
SET refresh_token = $1,
tkn_expires_at = $2
WHERE id = $3;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, password FROM users
WHERE email = $1; 

-- name: UpdateUserEmail :one
UPDATE users
SET email = $2, updated_at = $3
WHERE id = $1
RETURNING *;

-- name: UpdateUserName :one
UPDATE users
SET name = $2, updated_at = $3
WHERE id = $1
RETURNING *;

-- name: UpdateUserPassword :one
UPDATE users
SET password = $2, updated_at = $3
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

