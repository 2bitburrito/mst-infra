-- +goose Up
-- +goose StatementBegin
CREATE TYPE licence_type_enum AS ENUM (
    'trial',
    'paid',
    'beta'
);

ALTER TABLE licenses
ALTER COLUMN license_type TYPE licence_type_enum
USING license_type::licence_type_enum;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE licenses
ALTER COLUMN license_type TYPE TEXT;
-- +goose StatementEnd
