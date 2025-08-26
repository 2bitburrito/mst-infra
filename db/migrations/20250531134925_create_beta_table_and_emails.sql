-- +goose Up
-- +goose StatementBegin
CREATE TABLE beta_licences (
  email TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS beta_licences;
-- +goose StatementEnd
