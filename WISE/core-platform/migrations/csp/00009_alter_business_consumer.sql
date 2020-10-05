-- +goose Up
-- business table
ALTER TABLE business ADD COLUMN submitted timestamp with time zone DEFAULT NULL;
-- consumer  table
ALTER TABLE consumer ADD COLUMN submitted timestamp with time zone DEFAULT NULL;

-- +goose Down
-- business  table
ALTER TABLE business DROP COLUMN submitted;
-- consumer  table
ALTER TABLE consumer DROP COLUMN submitted;