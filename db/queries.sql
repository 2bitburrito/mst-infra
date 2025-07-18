-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserFromEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetLicence :one
SELECT * FROM licences
WHERE licence_key = $1 LIMIT 1;

-- name: GetAllLicencesFromUserID :many
SELECT * FROM licences
WHERE user_id = $1;

-- name: ChangeJTI :exec
UPDATE licences
SET jti = $1
WHERE licence_key = $2;

-- name: ChangeMachineIDAndJTI :exec
UPDATE licences
SET machine_id = $2, jti = $3
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

-- name: GetNameFromBetaList :one
SELECT name FROM beta_licences
WHERE email = $1;

-- name: AddTrialLicence :one
INSERT INTO licences (user_id, machine_id, licence_type, expiry)
VALUES ($1, $2, 'trial', NOW() + INTERVAL '14 days')
RETURNING licence_key, expiry;

-- name: AddBetaLicence :one
INSERT INTO licences (user_id, licence_type, expiry)
VALUES ($1, 'beta', NOW() + INTERVAL '60 days')
RETURNING licence_key, expiry;

-- name: SetBetaRowToSeen :exec
UPDATE beta_licences
SET seen = true
WHERE email = $1;


-- name: UnsetIsLatest :exec
UPDATE app_releases 
SET is_latest = FALSE 
WHERE platform = $1 
  AND architecture = $2 
  AND is_latest = TRUE;

-- name: AddNewReleaseData :exec
INSERT INTO app_releases (
  platform, 
  architecture,
  release_version,
  url_filename,
  file_size,
  release_date,
  is_latest,
  release_notes
  ) 
VALUES ($1, $2, $3, $4, $5, $6, TRUE, $7);

-- name: GetLatestBinary :one
SELECT * FROM app_releases
WHERE is_latest = TRUE
  AND architecture = $1
  AND platform = $2
LIMIT 1;
