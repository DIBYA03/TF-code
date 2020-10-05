-- +goose Up
DROP INDEX business_transaction_bank_transaction_id_bank_name_idx;
DROP INDEX business_pending_transaction_bank_transaction_id_bank_name_idx;
ALTER TABLE business_transaction ADD COLUMN notification_id UUID DEFAULT NULL;
CREATE UNIQUE INDEX business_transaction_notification_id_idx ON business_transaction (notification_id);
ALTER TABLE business_pending_transaction ADD COLUMN notification_id UUID DEFAULT NULL;
CREATE UNIQUE INDEX business_pending_transaction_notification_id_idx ON business_pending_transaction (notification_id);

-- +goose Down
DROP INDEX business_pending_transaction_notification_id_idx;
ALTER TABLE business_pending_transaction DROP COLUMN notification_id;
DROP INDEX business_transaction_notification_id_idx;
ALTER TABLE business_transaction DROP COLUMN notification_id;
CREATE UNIQUE INDEX business_pending_transaction_bank_transaction_id_bank_name_idx ON business_pending_transaction (bank_transaction_id, bank_name);
CREATE UNIQUE INDEX business_transaction_bank_transaction_id_bank_name_idx ON business_transaction (bank_transaction_id, bank_name);