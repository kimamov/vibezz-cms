-- name: GetUserByEmail :one
SELECT id, email, name, password_hash, role, created_at, updated_at
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT id, email, name, password_hash, role, created_at, updated_at
FROM users
WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (id, email, name, password_hash, role, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, email, name, password_hash, role, created_at, updated_at;

-- name: CountAdminUsers :one
SELECT COUNT(*) FROM users WHERE role = 'admin';
