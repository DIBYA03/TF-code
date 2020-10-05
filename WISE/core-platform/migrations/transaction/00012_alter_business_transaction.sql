-- +goose Up
/* PLATFORM-374: Add request_id to link to business_money_request table */
ALTER TABLE business_transaction ADD COLUMN money_request_id uuid DEFAULT NULL;

-- +goose Down
ALTER TABLE business_transaction DROP COLUMN money_request_id;
