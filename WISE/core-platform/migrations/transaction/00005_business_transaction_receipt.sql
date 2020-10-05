-- +goose Up
CREATE TABLE business_transaction_receipt ( 
    id                         UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(), 
    transaction_id             UUID REFERENCES business_transaction (id),
    created_user_id            UUID NOT NULL,
    business_id                UUID NOT NULL,
    content_type               TEXT DEFAULT NULL, 
    storage_key                TEXT DEFAULT NULL, 
    deleted                    TIMESTAMP WITH time zone DEFAULT NULL, 
    created                    TIMESTAMP WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP, 
    modified                   TIMESTAMP WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP 
);

CREATE TRIGGER update_business_transaction_receipt_modified BEFORE UPDATE ON business_transaction_receipt FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

CREATE INDEX business_transaction_receipt_transaction_id_fk ON business_transaction_receipt(transaction_id);
CREATE INDEX business_transaction_receipt_created_user_id_idx ON business_transaction_receipt(created_user_id);
CREATE INDEX business_transaction_receipt_business_id_idx ON business_transaction_receipt(business_id);


-- +goose Down
DROP TRIGGER IF EXISTS update_business_transaction_receipt_modified ON business_transaction_receipt;
DROP INDEX business_transaction_receipt_business_id_idx;
DROP INDEX business_transaction_receipt_created_user_id_idx;
DROP INDEX business_transaction_receipt_transaction_id_fk;
DROP TABLE business_transaction_receipt;
