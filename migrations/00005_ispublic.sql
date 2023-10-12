-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
ALTER TABLE galleries ADD COLUMN uuid TEXT;
UPDATE galleries SET uuid = substr(gen_random_uuid()::text, 1, 10);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE galleries DROP COLUMN uuid;
-- +goose StatementEnd
