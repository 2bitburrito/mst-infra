-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

ALTER TABLE licenses DROP CONSTRAINT IF EXISTS licenses_user_id_fkey;

ALTER TABLE users DROP COLUMN "id";
ALTER TABLE users ADD COLUMN "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY ;

ALTER TABLE licenses DROP COLUMN "user_id";
ALTER TABLE licenses ADD COLUMN "user_id" UUID DEFAULT gen_random_uuid();

ALTER TABLE licenses ADD CONSTRAINT licenses_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id);
-- +goose StatementEnd 

-- +goose Down
-- +goose StatementBegin
ALTER TABLE licenses ADD COLUMN "user_id" integer PRIMARY KEY;
ALTER TABLE users ADD COLUMN "id" integer;
-- +goose StatementEnd
