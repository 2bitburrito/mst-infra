-- +goose Up
-- +goose StatementBegin
ALTER TABLE licences
   ADD COLUMN jti UUID;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE licences
   DROP COLUMN jti UUID;
-- +goose StatementEnd
