-- +goose Up
-- +goose StatementBegin
ALTER TABLE licences
ADD COLUMN expiry DATE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE licences
DROP COLUMN expiry;
-- +goose StatementEnd
