-- +goose Up
/* Consumer document table */
CREATE TABLE consumer_document (
	id                              uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	consumer_id                         uuid NOT NULL REFERENCES consumer (id),
	number                          text DEFAULT NULL,
	doc_type                        text DEFAULT NULL,
	issuing_state                   text DEFAULT NULL,
	issuing_country                 text DEFAULT NULL,
	issued_date                     date DEFAULT NULL,
	expiration_date                 date DEFAULT NULL,
	content_type                    text NOT NULL,
	storage_key                     text DEFAULT NULL,
	deleted                         timestamp with time zone DEFAULT NULL,
    content_uploaded                timestamp with time zone DEFAULT NULL,
	created                         timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
	modified                        timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX consumer_document_consumer_id_fkey ON consumer_document (consumer_id);

CREATE TRIGGER update_consumer_document_modified BEFORE UPDATE ON consumer_document FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_consumer_document_modified on consumer_document;
DROP INDEX consumer_document_consumer_id_fkey;
DROP TABLE consumer_document;
