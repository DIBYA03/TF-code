
-- +goose Up
 CREATE TABLE consumer_document (
    id          uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    document_id uuid UNIQUE  NOT NULL,
    submitted   timestamp with time zone DEFAULT NULL,
    created     timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    modified     timestamp with time zone DEFAULT CURRENT_TIMESTAMP
 );
CREATE TRIGGER update_consumer_document_modified BEFORE UPDATE ON consumer_document FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
CREATE INDEX consumer_document_document_id_idx ON consumer_document(document_id);
CREATE INDEX consumer_document_submitted_idx ON consumer_document(submitted);

-- +goose Down
DROP TRIGGER IF EXISTS update_consumer_document_modified  ON consumer_document;
DROP INDEX consumer_document_document_id_idx;
DROP INDEX consumer_document_submitted_idx;
DROP TABLE consumer_document;
