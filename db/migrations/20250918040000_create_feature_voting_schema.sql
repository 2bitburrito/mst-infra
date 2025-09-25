-- +goose Up
-- +goose StatementBegin

CREATE TABLE features (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    created_by TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE votes (
    user_email TEXT NOT NULL,
    feature_id INTEGER NOT NULL,
    PRIMARY KEY (user_email),
    FOREIGN KEY (feature_id) REFERENCES features(id) ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS votes;
DROP TABLE IF EXISTS features;

-- +goose StatementEnd
