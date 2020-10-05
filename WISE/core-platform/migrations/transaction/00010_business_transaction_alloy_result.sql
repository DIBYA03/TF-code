-- +goose Up
/* Business transaction alloy result */
CREATE TABLE business_transaction_alloy_result ( 
	id                         UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(), 
	transaction_id             UUID NOT NULL REFERENCES business_transaction (id), 
	transaction_code_type      TEXT NOT NULL,
	transaction_type           TEXT DEFAULT NULL,
	amount                     DOUBLE PRECISION NOT NULL,
	currency                   TEXT NOT NULL,
	result                     TEXT NOT NULL,
	score                      DOUBLE PRECISION NOT NULL DEFAULT 0,
	outcome                    TEXT NOT NULL, 
	alloy_fraud_score          TEXT DEFAULT NULL,
	result_data                JSONB NOT NULL DEFAULT '{}'::jsonb, 
	created                    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp 
); 

CREATE INDEX business_transaction_alloy_result_transaction_id_fk ON business_transaction_alloy_result(transaction_id);

-- +goose Down
DROP INDEX business_transaction_alloy_result_transaction_id_fk;
DROP TABLE business_transaction_alloy_result;
