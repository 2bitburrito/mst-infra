-- +goose Up
-- +goose StatementBegin
ALTER TABLE beta_licences
   ADD COLUMN name TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE beta_licences
   DROP COLUMN name;
-- +goose StatementEnd
