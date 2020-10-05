-- +goose Up
ALTER TABLE business_transaction RENAME COLUMN notes TO source_notes;
ALTER TABLE business_pending_transaction RENAME COLUMN notes TO source_notes;
ALTER TABLE business_transaction_receipt RENAME TO business_transaction_attachment;

-- +goose Down
ALTER TABLE business_transaction_attachment RENAME TO business_transaction_receipt;
ALTER TABLE business_pending_transaction RENAME COLUMN source_notes TO notes;
ALTER TABLE business_transaction RENAME COLUMN source_notes TO notes;
