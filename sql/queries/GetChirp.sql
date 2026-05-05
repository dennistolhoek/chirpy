-- name: GetChirp :one
SELECT * FROM chirps
WHERE id = $1;

-- name: GetChirps :many
SELECT * FROM chirps
ORDER BY created_at ASC;