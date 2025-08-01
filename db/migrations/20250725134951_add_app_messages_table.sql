-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS app_messages (
    email text NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    number_of_licenses integer DEFAULT 0 NOT NULL,
    subscribed_to_emails boolean DEFAULT false NOT NULL,
    full_name character varying(255) DEFAULT ''::character varying NOT NULL,
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
