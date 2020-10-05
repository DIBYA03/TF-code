-- +goose Up
ALTER TABLE business ALTER COLUMN legal_name DROP NOT NULL;
ALTER TABLE business ALTER COLUMN legal_name SET DEFAULT NULL;

ALTER TABLE user_document DROP COLUMN issuing_auth;

ALTER TABLE user_document ALTER COLUMN number DROP NOT NULL;
ALTER TABLE user_document ALTER COLUMN number SET DEFAULT NULL;
ALTER TABLE user_document ALTER COLUMN doc_type DROP NOT NULL;
ALTER TABLE user_document ALTER COLUMN doc_type SET DEFAULT NULL;
ALTER TABLE user_document ALTER COLUMN content_type DROP DEFAULT;
ALTER TABLE user_document ALTER COLUMN content_type SET NOT NULL;

ALTER TABLE business_document DROP COLUMN issuing_auth;

ALTER TABLE business_document ALTER COLUMN number DROP NOT NULL;
ALTER TABLE business_document ALTER COLUMN number SET DEFAULT NULL;
ALTER TABLE business_document ALTER COLUMN doc_type DROP NOT NULL;
ALTER TABLE business_document ALTER COLUMN doc_type SET DEFAULT NULL;
ALTER TABLE business_document ALTER COLUMN content_type DROP DEFAULT;
ALTER TABLE business_document ALTER COLUMN content_type SET NOT NULL;

-- +goose Down
ALTER TABLE business_document ALTER COLUMN content_type DROP NOT NULL;
ALTER TABLE business_document ALTER COLUMN content_type SET DEFAULT NULL;
ALTER TABLE business_document ALTER COLUMN doc_type DROP DEFAULT;
ALTER TABLE business_document ALTER COLUMN doc_type SET NOT NULL;
ALTER TABLE business_document ALTER COLUMN number DROP DEFAULT;
ALTER TABLE business_document ALTER COLUMN number SET NOT NULL;

ALTER TABLE business_document ADD COLUMN issuing_auth text DEFAULT NULL;

ALTER TABLE user_document ALTER COLUMN content_type DROP NOT NULL;
ALTER TABLE user_document ALTER COLUMN content_type SET DEFAULT NULL;
ALTER TABLE user_document ALTER COLUMN doc_type DROP DEFAULT;
ALTER TABLE user_document ALTER COLUMN doc_type SET NOT NULL;
ALTER TABLE user_document ALTER COLUMN number DROP DEFAULT;
ALTER TABLE user_document ALTER COLUMN number SET NOT NULL;


ALTER TABLE user_document ADD COLUMN issuing_auth text DEFAULT NULL;

ALTER TABLE business ALTER COLUMN legal_name DROP DEFAULT;
ALTER TABLE business ALTER COLUMN legal_name SET NOT NULL;
