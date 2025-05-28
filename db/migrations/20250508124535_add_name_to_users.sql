-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN full_name VARCHAR(255) NOT NULL DEFAULT '';
ALTER TABLE users ALTER COLUMN number_of_licenses SET DEFAULT 0;
ALTER TABLE users ALTER COLUMN subscribed_to_emails SET DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users ALTER COLUMN subscribed_to_emails DROP DEFAULT;
ALTER TABLE users ALTER COLUMN number_of_licenses DROP DEFAULT;
ALTER TABLE users DROP COLUMN name;
--
-- +goose StatementEnd
