-- +goose Up
CREATE TABLE business_card_pending_transaction ( 
	id                          uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(), 
	transaction_id              uuid NOT NULL REFERENCES business_pending_transaction(id), 
	cardholder_id               text DEFAULT NULL,         
	card_transaction_id         text DEFAULT NULL, 
	transaction_network         text NOT NULL, 
	auth_amount                 DECIMAL(19,4) NOT NULL, 
	auth_date                   timestamp with time zone NOT NULL, 
	auth_response_code          text NOT NULL, 
	auth_number                 text DEFAULT NULL, 
	transaction_type            text NOT NULL, 
	local_amount                DECIMAL(19,4) NOT NULL, 
	local_currency              text NOT NULL, 
	local_date                  timestamp with time zone NOT NULL, 
	billing_currency            text NOT NULL,
	pos_entry_mode              text NOT NULL, 
	pos_condition_code          text NOT NULL, 
	acquirer_bin                text NOT NULL, 
	merchant_id                 text NOT NULL, 
	merchant_category_code      text NOT NULL, 
	merchant_terminal           text NOT NULL, 
	merchant_name               text NOT NULL, 
	merchant_street_address     text NOT NULL, 
	merchant_city               text NOT NULL, 
	merchant_state              text NOT NULL, 
	merchant_country            text NOT NULL, 
	created                     timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP 
); 

CREATE INDEX business_card_pending_transaction_transaction_id_fk ON business_card_pending_transaction(transaction_id);

-- +goose Down
DROP INDEX business_card_pending_transaction_transaction_id_fk;
DROP TABLE business_card_pending_transaction;
