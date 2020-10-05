-- +goose Up
ALTER TABLE business ADD COLUMN website TEXT DEFAULT NULL;
ALTER TABLE business ADD COLUMN online_info TEXT DEFAULT NULL;

-- +goose Down
ALTER TABLE business DROP COLUMN online_info;
ALTER TABLE business DROP COLUMN website;

