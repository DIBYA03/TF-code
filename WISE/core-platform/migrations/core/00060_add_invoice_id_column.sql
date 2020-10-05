-- +goose Up
ALTER TABLE business_money_request_payment ADD COLUMN invoice_id UUID;
-- +goose Down
ALTER TABLE business_money_request_payment DROP COLUMN invoice_id;