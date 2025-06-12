-- +goose Up
-- +goose StatementBegin
CREATE TABLE app_releases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    platform VARCHAR(20) NOT NULL, -- 'mac', 'windows', 'linux'
    architecture VARCHAR(10), -- 'x64', 'arm64', null for universal
    release_version VARCHAR(50) NOT NULL,
    url_filename VARCHAR(255) NOT NULL,
    file_size BIGINT,
    release_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_latest BOOLEAN DEFAULT FALSE,
    release_notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS app_releases;
-- +goose StatementEnd
