-- +goose Up
CREATE TABLE business_transaction_annotation ( 
    id                         UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(), 
    transaction_id             UUID REFERENCES business_transaction(id), 
    transaction_notes          text DEFAULT NULL,
    created                    TIMESTAMP WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP, 
    modified                   TIMESTAMP WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP 
);

CREATE INDEX business_transaction_annotation_transaction_id_fk ON business_card_transaction(transaction_id);

-- +goose Down
DROP INDEX business_transaction_annotation_transaction_id_fk;
DROP TABLE business_transaction_annotation;