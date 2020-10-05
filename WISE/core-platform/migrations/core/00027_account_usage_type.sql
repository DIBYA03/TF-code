-- +goose Up
ALTER TABLE business_bank_account ADD COLUMN usage_type TEXT NOT NULL DEFAULT 'primary';

ALTER TABLE business_linked_bank_account ADD COLUMN usage_type TEXT DEFAULT NULL;
UPDATE business_linked_bank_account SET usage_type = 'primary' WHERE business_bank_account_id IS NOT NULL;

-- +goose Down
ALTER TABLE business_linked_bank_account DROP COLUMN usage_type;

ALTER TABLE business_bank_account DROP COLUMN usage_type;
