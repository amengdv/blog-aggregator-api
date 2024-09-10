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
SELECT * FROM users
WHERE email = $1; 

-- name: UpdateUser :one
UPDATE users
SET email = $2, name = $3, password = $4, updated_at = $5
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

