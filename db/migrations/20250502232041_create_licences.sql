-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    email text NOT NULL,
    has_license boolean NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    number_of_licenses integer DEFAULT 0 NOT NULL,
    subscribed_to_emails boolean DEFAULT false NOT NULL,
    full_name character varying(255) DEFAULT ''::character varying NOT NULL,
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS licenses (
    license_key TEXT PRIMARY KEY,
    user_id uuid NOT NULL,
    machine_id TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_used_at TIMESTAMP NOT NULL,
    license_type TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS licenses;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
