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

-- name: InsertUser :exec
INSERT INTO users (id, email, full_name, has_license, subscribed_to_emails) 
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: UpdateUserId :exec
UPDATE users
SET id = $1
WHERE email = $2;

-- name: GetBetaEmail :one
SELECT * FROM beta_licences
WHERE email = $1;
