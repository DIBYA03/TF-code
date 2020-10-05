-- +goose Up
/* User document table */
CREATE TABLE user_document (
	id                              uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	user_id                         uuid NOT NULL REFERENCES wise_user (id),
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

CREATE INDEX user_document_user_id_fkey ON user_document (user_id);

CREATE TRIGGER update_user_document_modified BEFORE UPDATE ON user_document FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_user_document_modified on user_document;
DROP INDEX user_document_user_id_fkey;
DROP TABLE user_document;
