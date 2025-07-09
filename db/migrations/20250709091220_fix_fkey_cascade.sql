-- +goose Up
-- +goose StatementBegin
ALTER TABLE licences
DROP CONSTRAINT licences_user_id_fkey,
ADD CONSTRAINT licences_user_id_fkey
FOREIGN KEY (user_id) REFERENCES users(id)
ON UPDATE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE licences
DROP CONSTRAINT licences_user_id_fkey;
-- +goose StatementEnd
