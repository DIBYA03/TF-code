-- +goose Up
ALTER TABLE business_transaction_alloy_result DROP CONSTRAINT business_transaction_alloy_result_transaction_id_fkey;

DROP TABLE IF EXISTS business_transaction_alloy_result;

ALTER TABLE business_transaction_annotation DROP CONSTRAINT business_transaction_annotation_transaction_id_fkey;

ALTER TABLE business_transaction_attachment DROP CONSTRAINT business_transaction_receipt_transaction_id_fkey;

ALTER TABLE business_transaction_dispute DROP CONSTRAINT business_transaction_dispute_transaction_id_fkey;

-- +goose Down
ALTER TABLE business_transaction_dispute ADD CONSTRAINT business_transaction_dispute_transaction_id_fkey FOREIGN KEY (transaction_id) REFERENCES business_transaction (id);

ALTER TABLE business_transaction_attachment ADD CONSTRAINT business_transaction_receipt_transaction_id_fkey FOREIGN KEY (transaction_id) REFERENCES business_transaction (id);

ALTER TABLE business_transaction_annotation ADD CONSTRAINT business_transaction_annotation_transaction_id_fkey FOREIGN KEY (transaction_id) REFERENCES business_transaction (id);
