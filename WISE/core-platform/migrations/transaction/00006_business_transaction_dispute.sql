-- +goose Up
CREATE TABLE business_transaction_dispute ( 
    id                         UUID PRIMARY KEY NOT NULL DEFAULT Uuid_generate_v4(), 
    dispute_number             TEXT NOT NULL,
    transaction_id             UUID REFERENCES business_transaction (id), 
    receipt_id                 UUID DEFAULT NULL REFERENCES business_transaction_receipt (id), 
    created_user_id            UUID NOT NULL,
    business_id                UUID NOT NULL,
    category                   TEXT NOT NULL, 
    summary                    TEXT DEFAULT NULL, 
    dispute_status             TEXT DEFAULT NULL, 
    created                    TIMESTAMP WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP, 
    modified                   TIMESTAMP WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP 
); 

CREATE INDEX business_transaction_dispute_transaction_id_fk ON business_transaction_dispute(transaction_id);
CREATE INDEX business_transaction_dispute_receipt_id_fk ON business_transaction_dispute(receipt_id);
CREATE INDEX business_transaction_dispute_created_user_id_idx ON business_transaction_dispute(created_user_id);
CREATE INDEX business_transaction_dispute_business_id_idx ON business_transaction_dispute(business_id);

-- +goose Down
DROP INDEX business_transaction_dispute_business_id_idx;
DROP INDEX business_transaction_dispute_created_user_id_idx;
DROP INDEX business_transaction_dispute_receipt_id_fk;
DROP INDEX business_transaction_dispute_transaction_id_fk;
DROP TABLE business_transaction_dispute;
