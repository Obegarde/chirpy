
-- name: CreateChirp :one
INSERT INTO chirps(id,created_at,updated_at,body,user_id)
VALUES (
	gen_random_uuid(),
	NOW(),
	NOW(),
	$1,
	$2
	)
RETURNING *;

-- name: ResetChirps :exec
DELETE FROM chirps;

-- name: GetAllChirps :many
SELECT * 
FROM chirps
ORDER BY created_at ASC;

-- name: GetChirpByID :one
SELECT *
FROM chirps
WHERE ID = $1;