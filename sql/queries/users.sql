-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(), now(), now(), $1, $2
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET updated_at = now(), email = $2, hashed_password = $3
WHERE id = $1
RETURNING *;  

-- name: UpdateUserChirpyRed :exec
UPDATE users
SET updated_at = now(), is_chirpy_red = TRUE
WHERE id = $1;
