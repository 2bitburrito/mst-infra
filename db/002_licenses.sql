CREATE TABLE IF NOT EXISTS licenses (
    license_key TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    machine_id TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_used_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

