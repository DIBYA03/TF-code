-- +goose Up
/* PLATFORM-343: Add debit posted transaction ID and credit posted transaction ID  */
ALTER TABLE business_money_transfer RENAME COLUMN posted_transaction_id TO posted_debit_transaction_id;
ALTER TABLE business_money_transfer ADD COLUMN posted_credit_transaction_id UUID DEFAULT NULL;
CREATE INDEX business_money_transfer_posted_credit_transaction_id_idx ON business_money_transfer (posted_credit_transaction_id);

-- +goose Down
DROP INDEX business_money_transfer_posted_credit_transaction_id_idx;
ALTER TABLE business_money_transfer DROP COLUMN posted_credit_transaction_id;
ALTER TABLE business_money_transfer RENAME COLUMN posted_debit_transaction_id TO posted_transaction_id;
