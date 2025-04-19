-- name: GetUser :one
SELECT id, created_at, updated_at FROM users WHERE name = $1;
