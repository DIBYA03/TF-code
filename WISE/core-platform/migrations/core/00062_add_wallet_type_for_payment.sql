-- +goose Up
ALTER TABLE business_money_request_payment ADD COLUMN wallet_type text NULL;
-- +goose Down
ALTER TABLE business_money_request_payment DROP COLUMN wallet_type;