
-- +goose Up
ALTER TABLE consumer DROP COLUMN resolved;

-- +goose Down
ALTER TABLE consumer ADD COLUMN resolved timestamp with time zone DEFAULT NULL;