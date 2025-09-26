-- +goose Up
-- +goose StatementBegin
CREATE TABLE trial_machines (
  machine_id TEXT PRIMARY KEY,
  user_id UUID,
  claimed_at TIMESTAMP DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE trial_machines IF EXISTS;
-- +goose StatementEnd
