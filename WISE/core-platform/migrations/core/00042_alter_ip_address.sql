-- +goose Up
ALTER TABLE business_money_request ADD COLUMN request_ip_address TEXT DEFAULT NULL;

-- +goose Down
ALTER TABLE business_money_request DROP COLUMN request_ip_address;
