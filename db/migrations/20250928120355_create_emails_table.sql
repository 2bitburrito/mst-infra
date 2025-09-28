-- +goose Up
-- +goose StatementBegin
CREATE TYPE email_types AS ENUM ('marketing_invite', 'trial_licence_expiry');

CREATE TABLE sent_emails (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    user_id UUID REFERENCES users(id),
    email_type email_types NOT NULL,
    recipient_email TEXT NOT NULL,
    sent_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    status TEXT NOT NULL CHECK (status IN ('sent', 'failed', 'queued')),
    error_message TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sent_emails;
DROP TYPE IF EXISTS email_types;
-- +goose StatementEnd
