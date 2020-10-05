-- +goose Up
/* Add UUID extension  */
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- +goose StatementBegin
/* Add modified update function */
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.modified = now();
    RETURN NEW;
END;
$$ language 'plpgsql';
-- +goose StatementEnd

/* Business bank transaction table */
CREATE TABLE business_transaction (
	id                          uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	business_id                 uuid NOT NULL,
	bank_name                   text NOT NULL,
	bank_transaction_id         text NOT NULL,
	bank_extra                  jsonb NOT NULL DEFAULT '{}'::jsonb,
	transaction_type            text DEFAULT NULL,
	account_id                  uuid DEFAULT NULL,
	card_id                     uuid DEFAULT NULL,
	code_type                   text NOT NULL,
	amount                      double PRECISION NOT NULL,
	currency                    text NOT NULL,
	money_transfer_id           uuid DEFAULT NULL,
	contact_id                  uuid DEFAULT NULL,
	money_transfer_desc         text DEFAULT NULL,
	transaction_desc            text NOT NULL,
	transaction_date            timestamp with time zone NOT NULL,
	created                     timestamp with time zone NOT NULL DEFAULT current_timestamp
);

CREATE INDEX business_transaction_business_id_idx ON business_transaction(business_id);
CREATE INDEX business_transaction_business_id_amount_idx ON business_transaction (business_id, amount);
CREATE INDEX business_transaction_business_id_code_type_idx ON business_transaction (business_id, code_type);
CREATE INDEX business_transaction_business_id_transaction_date_idx ON business_transaction (business_id, transaction_date);
CREATE UNIQUE INDEX business_transaction_bank_transaction_id_bank_name_idx ON business_transaction (bank_transaction_id, bank_name);

-- +goose Down
DROP INDEX business_transaction_bank_transaction_id_bank_name_idx;
DROP INDEX business_transaction_business_id_transaction_date_idx;
DROP INDEX business_transaction_business_id_code_type_idx;
DROP INDEX business_transaction_business_id_amount_idx;
DROP INDEX business_transaction_business_id_idx;
DROP TABLE business_transaction;

DROP FUNCTION IF EXISTS update_modified_column;
