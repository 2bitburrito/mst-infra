-- +goose Up
-- +goose StatementBegin
ALTER TABLE beta_licences
  ADD COLUMN seen BOOLEAN NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE beta_licences
  DROP COLUMN seen;
-- +goose StatementEnd
