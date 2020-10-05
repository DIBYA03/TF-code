
-- +goose Up
CREATE TABLE business_hold_pending_transaction ( 
	id                          uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(), 
	transaction_id              uuid NOT NULL REFERENCES business_pending_transaction(id), 
	amount                      DECIMAL(19,4) NOT NULL, 
	hold_number                 text NOT NULL, 
	transaction_date            timestamp with time zone NOT NULL, 
	expiry_date                 timestamp with time zone NOT NULL, 
	created                     timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP 
); 

CREATE INDEX business_hold_pending_transaction_transaction_id_fk ON business_hold_pending_transaction(transaction_id);

-- +goose Down
DROP INDEX business_hold_pending_transaction_transaction_id_fk;
DROP TABLE business_hold_pending_transaction;
