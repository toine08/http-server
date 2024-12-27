-- name: CreateUser :one 
INSERT INTO users(id, created_at, updated_at, email)
VALUES(
gen_random_uuid(),
TIMESTAMP.NOW(),
TIMESTAMP.NOW(),
$1
)
RETURNING *;