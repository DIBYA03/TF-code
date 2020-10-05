-- +goose Up
ALTER TABLE business_transaction ADD COLUMN notes TEXT DEFAULT NULL;

-- +goose Down
ALTER TABLE business_transaction DROP COLUMN notes;
