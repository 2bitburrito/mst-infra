-- +goose Up
-- +goose StatementBegin

CREATE TABLE stripe_events (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL,
    processed_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE if EXISTS stripe_events;
-- +goose StatementEnd
