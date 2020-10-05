-- +goose Up
ALTER TABLE business_money_request_payment ADD COLUMN linked_bank_account_id uuid REFERENCES business_linked_bank_account(id) DEFAULT NULL;

-- +goose Down
ALTER TABLE business_money_request_payment DROP COLUMN linked_bank_account_id;
