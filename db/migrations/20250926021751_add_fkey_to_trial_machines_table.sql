-- +goose Up
-- +goose StatementBegin
ALTER TABLE trial_machines
ADD CONSTRAINT trial_machines_user_id_fkey
FOREIGN KEY (user_id) REFERENCES users(id)
ON UPDATE CASCADE
ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE licences
DROP CONSTRAINT trial_machines_user_id_fkey;
-- +goose StatementEnd

