-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
   DROP has_license;
ALTER TABLE users
  RENAME COLUMN number_of_licenses TO number_of_licences;   

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
  ADD COLUMN has_license BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE users
  RENAME COLUMN number_of_licences TO number_of_licenses;

-- +goose StatementEnd
