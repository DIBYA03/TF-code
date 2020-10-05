-- +goose Up
/* Signature requests table */
CREATE TABLE signature_request (
    id                          uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    business_id                 uuid REFERENCES business(id),
    template_type               text NOT NULL,
    template_provider           text NOT NULL,
    signature_request_id        text NOT NULL,
    signature_id                text NOT NULL,
    signature_status            text NOT NULL,
    document_id                 uuid REFERENCES business_document(id),
    created                     timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified                    timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX signature_request_business_id_fk ON signature_request(business_id);
CREATE INDEX signature_request_document_id_fk ON signature_request(document_id);
CREATE UNIQUE INDEX signature_request_business_id_template_type_idx
ON signature_request (business_id, template_type);

CREATE TRIGGER update_signature_request_modified BEFORE UPDATE ON signature_request FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_signature_request_modified on signature_request;

DROP INDEX signature_request_business_id_template_type_idx;
DROP INDEX signature_request_document_id_fk;
DROP INDEX signature_request_business_id_fk;
DROP TABLE signature_request;
