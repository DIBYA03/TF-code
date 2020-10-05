-- +goose Up
ALTER TABLE money_transfer_request ADD COLUMN expiration_date timestamp with time zone DEFAULT NULL;

-- +goose Down
ALTER TABLE money_transfer_request DROP COLUMN expiration_date;