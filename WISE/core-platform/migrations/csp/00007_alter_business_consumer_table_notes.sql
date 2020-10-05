
-- +goose Up
ALTER TABLE business ALTER COLUMN notes SET DEFAULT '[]'::jsonb;
ALTER TABLE business ALTER COLUMN notes TYPE jsonb USING notes::jsonb;

ALTER TABLE consumer ALTER COLUMN notes SET DEFAULT '[]'::jsonb;
ALTER TABLE consumer ALTER COLUMN notes TYPE jsonb USING notes::jsonb;

-- +goose Down
ALTER TABLE business ALTER COLUMN notes TYPE TEXT USING notes::TEXT;
ALTER TABLE business ALTER COLUMN notes SET DEFAULT ''::TEXT;

ALTER TABLE consumer ALTER COLUMN notes TYPE TEXT USING notes::TEXT;
ALTER TABLE consumer ALTER COLUMN notes SET DEFAULT ''::TEXT;
