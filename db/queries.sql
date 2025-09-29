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
INSERT INTO users (id, email, full_name, subscribed_to_emails) 
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: UpdateUserId :exec
UPDATE users
SET id = $1
WHERE email = $2;

-- name: GetBetaEmail :one
SELECT * FROM beta_licences
WHERE email = $1;

-- name: GetAllBetaEmails :many
SELECT * FROM beta_licences;

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

-- name: UpdateLicenceUsedTime :exec
UPDATE licences
SET last_used_at = NOW()
WHERE licence_key = $1;


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

-- name: AddStripeIDtoUser :exec
UPDATE users
SET stripe_id = $2
WHERE id = $1;

-- name: AddStripeEvent :exec
INSERT INTO stripe_events (
  id,
  type,
  processed_at
  )
  VALUES($1, $2, NOW());

-- name: StripeEventExists :one
SELECT EXISTS(
  SELECT * FROM stripe_events WHERE id=$1
);

-- name: AddPaidLicence :one
INSERT INTO licences (
  user_id,
  created_at,
  licence_type
) 
VALUES (
  $1, NOW(), 'paid'
  ) RETURNING licence_key;

-- name: IncrementLicences :exec
UPDATE users
SET number_of_licences = number_of_licences + $2
WHERE id = $1;

-- name: MachineIDIsUsed :one
SELECT EXISTS(
    SELECT * FROM trial_machines
    WHERE machine_id = $1
    AND claimed_at <= now() - INTERVAL '14 days'
);

-- name: AddMachineID :exec
INSERT INTO trial_machines(
  machine_id, user_id
)VALUES(
  $1, $2
);

-- name: HasParticipatedInBeta :one
SELECT EXISTS(
  SELECT * FROM beta_licences
  WHERE email = $1 AND seen = 'TRUE'
);

-- name: GetExpiredLicences :many
SELECT l.*, u.* FROM licences l
  JOIN users u ON l.user_id = u.id
  WHERE l.expiry < now();

-- name: GetSentTrialLicenceExpiryEmails :many
SELECT * FROM sent_emails
  WHERE user_id = $1
  AND email_type = 'trial_licence_expiry';

-- name: SendEmail :exec
INSERT INTO sent_emails(
  user_id,
  email_type,
  recipient_email,
  status,
  error_message
) VALUES ($1, $2, $3, $4, $5);
