-- +goose Up
ALTER TABLE business_money_request ADD COLUMN request_source TEXT DEFAULT NULL;
ALTER TABLE business_money_request ADD COLUMN request_source_id TEXT DEFAULT NULL;

-- +goose Down
ALTER TABLE business_money_request DROP COLUMN request_source_id;
ALTER TABLE business_money_request DROP COLUMN request_source;