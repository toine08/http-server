-- name: CreateUser :one 
INSERT INTO users(id, created_at, updated_at, email,hashed_password)
VALUES(
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;
-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, hashed_password FROM users WHERE email =$1;

-- name: GetUserById :one
SELECT id, created_at, updated_at, email, hashed_password FROM users WHERE id =$1;

-- name: UpdateUser :one
UPDATE users
SET email = $2, hashed_password = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;



