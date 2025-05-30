-- +goose Up
-- +goose StatementBegin
-- ALTER TABLE users ADD COLUMN number_of_licenses INTEGER NOT NULL;
-- ALTER TABLE users ADD COLUMN subscribed_to_emails BOOLEAN NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
