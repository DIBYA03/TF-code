-- +goose Up
/* Add UUID extension  */
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

/* Business bank pending transaction table */
CREATE TABLE business_pending_transaction (
	id                          uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	business_id                 uuid NOT NULL,
	bank_name                   text NOT NULL,
	bank_transaction_id         text DEFAULT NULL,
	bank_extra                  jsonb NOT NULL DEFAULT '{}'::jsonb,
	transaction_type            text DEFAULT NULL,
	account_id                  uuid DEFAULT NULL,
	card_id                     uuid DEFAULT NULL,
	code_type                   text NOT NULL,
	amount                      DECIMAL(19,4) NOT NULL,
	currency                    text NOT NULL,
	money_transfer_id           uuid DEFAULT NULL,
	contact_id                  uuid DEFAULT NULL,
	money_transfer_desc         text DEFAULT NULL,
	transaction_desc            text NOT NULL,
	transaction_date            timestamp with time zone NOT NULL,
	transaction_status          text DEFAULT NULL,
	partner_name                text NOT NULL,
	notes                       text DEFAULT NULL,
	money_request_id            uuid DEFAULT NULL,
	created                     timestamp with time zone NOT NULL DEFAULT current_timestamp
);


CREATE INDEX business_pending_transaction_business_id_idx ON business_pending_transaction(business_id);
CREATE INDEX business_pending_transaction_business_id_amount_idx ON business_pending_transaction (business_id, amount);
CREATE INDEX business_pending_transaction_business_id_code_type_idx ON business_pending_transaction (business_id, code_type);
CREATE INDEX business_pending_transaction_business_id_transaction_date_idx ON business_pending_transaction (business_id, transaction_date);
CREATE UNIQUE INDEX business_pending_transaction_bank_transaction_id_bank_name_idx ON business_pending_transaction (bank_transaction_id, bank_name);

-- +goose Down
DROP INDEX business_pending_transaction_bank_transaction_id_bank_name_idx;
DROP INDEX business_pending_transaction_business_id_transaction_date_idx;
DROP INDEX business_pending_transaction_business_id_code_type_idx;
DROP INDEX business_pending_transaction_business_id_amount_idx;
DROP INDEX business_pending_transaction_business_id_idx;
DROP TABLE business_pending_transaction;
