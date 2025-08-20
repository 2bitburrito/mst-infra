-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD stripe_id TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN stripe_id;
-- +goose StatementEnd
