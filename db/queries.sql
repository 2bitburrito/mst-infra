-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetLicence :one
SELECT * FROM licences
where licence_key = $1 LIMIT 1;

-- name: ChangeMachineID :exec
UPDATE licences
SET machine_id = $2
WHERE licence_key = $1;

-- name: RemoveMachineID :exec
UPDATE licences
SET machine_id = null
WHERE licence_key = $1;
