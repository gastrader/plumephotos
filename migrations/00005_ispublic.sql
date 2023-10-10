-- +goose Up
-- +goose StatementBegin
ALTER TABLE galleries
ADD COLUMN is_public BOOLEAN DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE galleries DROP COLUMN is_public;
-- +goose StatementEnd
