-- +goose Up
-- +goose StatementBegin
ALTER TABLE licences
ADD CONSTRAINT licences_user_id_fkey
FOREIGN KEY (user_id) REFERENCES users(id)
ON UPDATE CASCADE
ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE licences
DROP CONSTRAINT licences_user_id_fkey;
-- +goose StatementEnd
