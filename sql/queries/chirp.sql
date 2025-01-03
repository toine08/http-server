-- name: CreateChirp :one 
INSERT INTO chirps(id, created_at, updated_at, body, user_id)
VALUES(
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteChirps :exec
DELETE FROM chirps;

-- name: AllChirps :many
SELECT * FROM chirps
ORDER BY created_at ASC;

-- name: ChirpByID :one
SELECT id, created_at, updated_at, body, user_id FROM chirps WHERE id = $1;

-- name: DeleteChirpByID :one
DELETE FROM chirps WHERE id = $1 AND user_id = $2
RETURNING *;