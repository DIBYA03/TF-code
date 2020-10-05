
-- +goose Up
-- business document table
ALTER TABLE business_document ADD COLUMN document_status text DEFAULT 'notStarted';
ALTER TABLE business_document ADD COLUMN response jsonb DEFAULT '{}'::jsonb;
-- consumer document table
ALTER TABLE consumer_document ADD COLUMN document_status text DEFAULT 'notStarted';
ALTER TABLE consumer_document ADD COLUMN response jsonb DEFAULT '{}'::jsonb;

-- +goose Down
-- business document table
ALTER TABLE business_document DROP COLUMN document_status;
ALTER TABLE business_document DROP COLUMN response;
-- consumer document table
ALTER TABLE consumer_document DROP COLUMN document_status;
ALTER TABLE consumer_document DROP COLUMN response;
