-- +goose Up
ALTER TABLE business ADD COLUMN logo_storage_key TEXT DEFAULT NULL;

-- +goose Down
ALTER TABLE business DROP COLUMN logo_storage_key;