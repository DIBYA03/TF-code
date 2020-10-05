-- +goose Up
ALTER TABLE business_transaction ADD COLUMN transaction_title TEXT DEFAULT NULL;
ALTER TABLE business_transaction ADD COLUMN transaction_subtype TEXT DEFAULT NULL;
ALTER TABLE business_transaction RENAME COLUMN money_transfer_desc TO bank_transaction_desc;

ALTER TABLE business_pending_transaction ADD COLUMN transaction_title TEXT DEFAULT NULL;
ALTER TABLE business_pending_transaction ADD COLUMN transaction_subtype TEXT DEFAULT NULL;
ALTER TABLE business_pending_transaction RENAME COLUMN money_transfer_desc TO bank_transaction_desc;

-- +goose Down
ALTER TABLE business_pending_transaction RENAME COLUMN bank_transaction_desc TO money_transfer_desc;
ALTER TABLE business_pending_transaction DROP COLUMN transaction_subtype;
ALTER TABLE business_pending_transaction DROP COLUMN transaction_title;

ALTER TABLE business_transaction RENAME COLUMN bank_transaction_desc TO money_transfer_desc;
ALTER TABLE business_transaction DROP COLUMN transaction_subtype;
ALTER TABLE business_transaction DROP COLUMN transaction_title;
