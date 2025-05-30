-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS licenses;


CREATE TYPE licence_type_enum AS ENUM (
    'trial',
    'paid',
    'beta'
);
CREATE TABLE IF NOT EXISTS licences (
    licence_key TEXT PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    machine_id TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_used_at TIMESTAMP,
    licence_type licence_type_enum,
    FOREIGN KEY (user_id) REFERENCES users (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS licences;
-- +goose StatementEnd
