-- name: GetChirp :one
SELECT * FROM chirps
WHERE id = $1;

-- name: GetChirps :many
SELECT * FROM chirps
ORDER BY created_at ASC;

-- name: GetChirpsByAuthor :many
SELECT * FROM chirps
WHERE User_ID = $1
ORDER BY created_at ASC;