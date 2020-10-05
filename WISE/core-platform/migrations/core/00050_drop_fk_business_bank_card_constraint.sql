-- +goose Up
ALTER TABLE business_bank_card DROP CONSTRAINT business_bank_card_bank_account_id_fkey;

-- +goose Down
ALTER TABLE business_bank_card ADD CONSTRAINT business_bank_card_bank_account_id_fkey FOREIGN KEY (bank_account_id) REFERENCES business_bank_account (id);
