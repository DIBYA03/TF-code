-- +goose Up
/* Business document table */
CREATE TABLE business_document (
	id                              uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	business_id                     uuid NOT NULL REFERENCES business (id),
	created_user_id                 uuid NOT NULL REFERENCES wise_user (id),
	number                          text NOT NULL,
	doc_type                        text NOT NULL,
	issuing_auth                    text DEFAULT NULL,
	issuing_state                   text DEFAULT NULL,
	issuing_country                 text DEFAULT NULL,
	issued_date                     date DEFAULT NULL,
	expiration_date                 date DEFAULT NULL,
	content_type                    text DEFAULT NULL,
	storage_key                     text DEFAULT NULL,
	deleted                         timestamp with time zone DEFAULT NULL,
	created                         timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
	modified                        timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX business_document_business_id_fkey ON business_document (business_id);
CREATE INDEX business_document_created_user_id_fkey ON business_document (created_user_id);

CREATE TRIGGER update_business_document_modified BEFORE UPDATE ON business_document FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_business_document_modified on business_document;
DROP INDEX business_document_created_user_id_fkey;
DROP INDEX business_document_business_id_fkey;
DROP TABLE business_document;
