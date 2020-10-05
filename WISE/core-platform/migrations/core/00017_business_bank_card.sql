-- +goose Up
/* Business bank card table */
CREATE TABLE business_bank_card ( 
    id                      UUID PRIMARY KEY NOT NULL DEFAULT Uuid_generate_v4(), 
    cardholder_id           UUID NOT NULL REFERENCES wise_user (id), 
    business_id             UUID NOT NULL REFERENCES business (id), 
    bank_account_id         UUID NOT NULL REFERENCES business_bank_account(id), 
    card_type               TEXT NOT NULL, 
    cardholder_name         TEXT NOT NULL, 
    is_virtual              BOOLEAN NOT NULL DEFAULT FALSE, 
    bank_name               TEXT NOT NULL, 
    bank_card_id            TEXT NOT NULL, 
    bank_extra              JSONB NOT NULL DEFAULT '{}'::jsonb, 
    card_number_masked      TEXT NOT NULL,
    card_brand              TEXT DEFAULT NULL, 
    currency                TEXT NOT NULL, 
    card_status             TEXT NOT NULL, 
    alias                   TEXT DEFAULT NULL, 
    daily_withdrawal_limit  DOUBLE PRECISION NOT NULL, 
    daily_pos_limit         DOUBLE PRECISION NOT NULL, 
    daily_transaction_limit INTEGER DEFAULT 0, 
    card_block              JSONB DEFAULT NULL,
    created                 TIMESTAMP WITH TIME zone NOT NULL DEFAULT current_timestamp, 
    modified                TIMESTAMP WITH TIME zone NOT NULL DEFAULT current_timestamp 
  ); 

CREATE INDEX business_bank_card_cardholder_id_fkey ON business_bank_card (cardholder_id);
CREATE INDEX business_bank_card_business_id_fkey ON business_bank_card (business_id);
CREATE INDEX business_bank_card_bank_account_id_fkey ON business_bank_card (bank_account_id);
CREATE UNIQUE INDEX business_bank_card_bank_card_id_bank_name_idx ON business_bank_card(bank_card_id, bank_name);

CREATE TRIGGER update_business_bank_card_modified BEFORE UPDATE ON business_bank_card FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_business_bank_card_modified on business_bank_card;
DROP INDEX business_bank_card_bank_card_id_bank_name_idx;
DROP INDEX business_bank_card_bank_account_id_fkey;
DROP INDEX business_bank_card_business_id_fkey;
DROP INDEX business_bank_card_cardholder_id_fkey;
DROP TABLE business_bank_card;
