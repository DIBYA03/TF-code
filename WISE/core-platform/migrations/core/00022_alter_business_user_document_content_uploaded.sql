
-- +goose Up
ALTER TABLE business_document ADD COLUMN content_uploaded timestamp with time zone DEFAULT NULL;
ALTER TABLE user_document ADD COLUMN content_uploaded timestamp with time zone DEFAULT NULL;

-- +goose Down
ALTER TABLE business_document DROP COLUMN content_uploaded;
ALTER TABLE user_document DROP COLUMN content_uploaded;
