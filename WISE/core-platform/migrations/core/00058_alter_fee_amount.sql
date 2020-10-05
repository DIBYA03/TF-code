-- +goose Up
ALTER TABLE business_money_request_payment ADD COLUMN fee_amount DECIMAL(19,2) NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE business_money_request_payment DROP COLUMN fee_amount;