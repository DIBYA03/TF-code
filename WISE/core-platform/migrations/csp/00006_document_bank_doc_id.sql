-- +goose Up
-- business document table
ALTER TABLE business_document ADD COLUMN bank_document_id text;
CREATE INDEX business_document_bank_document_id_idx ON business_document(bank_document_id);

-- consumer document table
ALTER TABLE consumer_document ADD COLUMN bank_document_id text;
CREATE INDEX consumer_document_bank_document_id_idx ON consumer_document(bank_document_id);

-- +goose Down
-- business document table
DROP INDEX business_document_bank_document_id_idx;
ALTER TABLE business_document DROP COLUMN bank_document_id;

-- consumer document table
DROP INDEX consumer_document_bank_document_id_idx;
ALTER TABLE consumer_document DROP COLUMN bank_document_id;
