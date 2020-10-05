 
-- +goose Up

 CREATE TABLE business_document (
    id          uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    document_id uuid UNIQUE  NOT NULL,
    submitted   timestamp with time zone DEFAULT NULL,
    created     timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    modified     timestamp with time zone DEFAULT CURRENT_TIMESTAMP
 );
CREATE TRIGGER update_business_document_modified BEFORE UPDATE ON business_document FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
CREATE INDEX business_document_document_id_idx ON business_document(document_id);
CREATE INDEX business_document_submitted_idx ON business_document(submitted);

-- +goose Down
DROP TRIGGER IF EXISTS update_business_document_modified  on business_document;
DROP INDEX business_document_submitted_idx;
DROP INDEX business_document_document_id_idx;
DROP TABLE business_document;
